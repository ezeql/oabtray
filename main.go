package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/gofrs/flock"
)

const (
	INITIAL_DISPLAY_DURATION  = 5 * time.Second
	SCREEN_WIDTH              = 20
	ANIMATION_SPEED           = 100 * time.Millisecond
	UPDATE_INTERVAL           = 30 * time.Second
	DATA_FILE                 = "bitcoin_tracker_data.gob"
	BULL_ANIMATION_DURATION   = time.Second
	BULL_ANIMATION_SPEED      = 100 * time.Millisecond
	LOCK_FILE                 = "bitcoin_tracker.lock"
	ALABA_FACTOR              = 2.5
)

var(
	VERSION = "dev"
)

type CoinGeckoResponse struct {
	Bitcoin struct {
		USD          float64 `json:"usd"`
		USD24HChange float64 `json:"usd_24h_change"`
	} `json:"bitcoin"`
}

type PersistentData struct {
	LastPrice         float64
	LastChangePercent float64
	LastUpdateTime    time.Time
}

type BinanceResponse struct {
	LastPrice      string  `json:"lastPrice"`
	PriceChange    string  `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
}

var (
	isFirstUpdate         = true
	mu                    sync.Mutex
	lastPrice             float64
	lastChangePercent     float64
	lastUpdateTime        time.Time
	isAnimating           bool = false
	animationMutex        sync.Mutex
	fileLock              *flock.Flock
	// lastAlabadoTime       time.Time
	toTheMoonMode         bool = false
)

func main() {
	fileLock = flock.New(getLockFilePath())
	locked, err := fileLock.TryLock()
	if err != nil {
		log.Printf("Error acquiring lock: %v", err)
		return
	}
	if !locked {
		fmt.Println("Another instance of the application is already running.")
		return
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			log.Printf("Error releasing lock: %v", err)
		}
	}()

	data := loadPersistentData()
	lastPrice = data.LastPrice
	lastChangePercent = data.LastChangePercent
	lastUpdateTime = data.LastUpdateTime

	go priceUpdater()
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("₿")
	systray.SetTooltip("Bitcoin Price Tracker")
	mBitcoinPrice := systray.AddMenuItem("OAB", "Bitcoin price")
	mBitcoinPrice.Disable()

	systray.AddSeparator()
	mMoonMode := systray.AddMenuItem("Set price in millions", "Toggle price in millions")
	if toTheMoonMode {
		mMoonMode.Check()
	}

	systray.AddSeparator()
	mVersion := systray.AddMenuItem("Version: " + VERSION, "Version information")
	mVersion.Disable()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mMoonMode.ClickedCh:
				toTheMoonMode = !toTheMoonMode
				if toTheMoonMode {
					mMoonMode.Check()
				} else {
					mMoonMode.Uncheck()
				}
				// Update display with current price
				if lastPrice > 0 {
					updateTrayQuiet(lastPrice, lastChangePercent)
				}
			}
		}
	}()

	if lastPrice > 0 {
		updateTrayQuiet(lastPrice, lastChangePercent)
		log.Printf("Loaded price from disk: $%.2f (%+.2f%%)", lastPrice, lastChangePercent)
	}
}

func onExit() {
	savePersistentData()
}

func priceUpdater() {
	time.Sleep(1 * time.Second)

	if time.Since(lastUpdateTime) > UPDATE_INTERVAL || lastUpdateTime.IsZero() {
		fetchAndUpdatePrice()
	}

	ticker := time.NewTicker(UPDATE_INTERVAL)
	defer ticker.Stop()

	for range ticker.C {
		fetchAndUpdatePrice()
	}
}

func fetchAndUpdatePrice() {
	price, changePercent, err := fetchPrice()
	if err != nil {
		log.Printf("Error fetching price: %v", err)
		displayError(err)
		return
	}

	mu.Lock()
	currentIsFirstUpdate := isFirstUpdate
	isFirstUpdate = false
	mu.Unlock()

	if currentIsFirstUpdate {
		updateTrayWithInitialDisplay(price, changePercent)
	} else {
		updateTray(price, changePercent)
	}

	lastPrice = price
	lastChangePercent = changePercent
	lastUpdateTime = time.Now()
	savePersistentData()
}

func fetchPrice() (float64, float64, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/24hr?symbol=BTCUSDT")
	if err != nil {
		return 0, 0, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	var data BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("JSON decode error: %v", err)
	}

	price, err := strconv.ParseFloat(data.LastPrice, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("price parse error: %v", err)
	}

	changePercent, err := strconv.ParseFloat(data.PriceChangePercent, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("change percent parse error: %v", err)
	}

	return price, changePercent, nil
}

func updateTrayWithInitialDisplay(price, changePercent float64) {
	priceStr := formatPriceString(price, changePercent)
	systray.SetTitle(priceStr)
	systray.SetTooltip(priceStr)

	time.Sleep(INITIAL_DISPLAY_DURATION)

	updateTray(price, changePercent)
}

func updateTray(price, changePercent float64) {
	updateTrayQuiet(price, changePercent)
	
	if math.Abs(changePercent) >= 5.0 {
		var animationText string
		if changePercent >= 5.0 {
			animationText = "ALABADO!!!"
		} else if changePercent <= -5.0 {
			animationText = "PUTA MADRE!"
		}
		
		if animationText != "" {
			go runAnimation(animationText)
		}
	}
}

func updateTrayQuiet(price, changePercent float64) {
	priceStr := formatPriceString(price, changePercent)
	systray.SetTitle(priceStr)
	systray.SetTooltip(priceStr)
}

func formatPriceString(price, changePercent float64) string {
	var emoji string
	if changePercent > 0 {
		emoji = "🟢"
	} else if changePercent < 0 {
		emoji = "🔴"
	} else {
		emoji = "⚪"
	}
	
	var priceStr string
	if toTheMoonMode {
		millions := price / 1000
		priceStr = fmt.Sprintf("%.3fM", millions/1000)
	} else {
		// Format with thousands separator using string manipulation
		priceStr = addThousandsSeparator(fmt.Sprintf("%.2f", price))
	}
	
	emoticons := getEmoticons(changePercent)
	return fmt.Sprintf("₿ %s $%s (%+.2f%%) %s", emoji, priceStr, changePercent, emoticons)
}

func addThousandsSeparator(s string) string {
	// Split the string into integer and decimal parts
	parts := strings.Split(s, ".")
	intPart := parts[0]
	
	// Add commas to the integer part
	var result []byte
	for i := len(intPart) - 1; i >= 0; i-- {
		if len(result) > 0 && (len(intPart)-i-1)%3 == 0 {
			result = append([]byte{','}, result...)
		}
		result = append([]byte{intPart[i]}, result...)
	}
	
	// Reconstruct the number with decimal part
	if len(parts) > 1 {
		return string(result) + "." + parts[1]
	}
	return string(result)
}

func getEmoticons(changePercent float64) string {
	absChangePercent := math.Abs(changePercent)
	count := int(math.Floor(absChangePercent / ALABA_FACTOR))
	if changePercent >= 0 {
		return strings.Repeat("🚀", count)
	} else {
		return strings.Repeat("🧂", count)
	}
}

func displayError(err error) {
	errorMsg := fmt.Sprintf("%-*s", SCREEN_WIDTH, "TECHNICAL DIFFICULTIES :)")
	systray.SetTitle(errorMsg)
	systray.SetTooltip(err.Error())
}

func runAnimation(text string) {
	animationMutex.Lock()
	if isAnimating {
		animationMutex.Unlock()
		return
	}
	isAnimating = true
	animationMutex.Unlock()

	// Show full text instantly
	paddedText := fmt.Sprintf("%-*s", SCREEN_WIDTH, text)
	systray.SetTitle(paddedText)
	time.Sleep(1 * time.Second)

	// Convert text to rune slice for character manipulation
	chars := []rune(text)
	positions := make([]int, len(chars))
	for i := range positions {
		positions[i] = i
	}

	// Randomly remove characters one by one
	for len(positions) > 0 {
		// Pick a random position to remove
		randIndex := rand.Intn(len(positions))
		removePos := positions[randIndex]
		
		// Remove the position from our slice
		positions = append(positions[:randIndex], positions[randIndex+1:]...)
		
		// Create display text with character removed
		displayChars := make([]rune, len(chars))
		copy(displayChars, chars)
		displayChars[removePos] = ' '
		
		// Update remaining characters
		for _, pos := range positions {
			displayChars[pos] = chars[pos]
		}
		
		// Display the result
		displayText := fmt.Sprintf("%-*s", SCREEN_WIDTH, string(displayChars))
		systray.SetTitle(displayText)
		time.Sleep(100 * time.Millisecond)
	}

	// Return to price display
	systray.SetTitle(formatPriceString(lastPrice, lastChangePercent))

	animationMutex.Lock()
	isAnimating = false
	animationMutex.Unlock()
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

func savePersistentData() {
	data := PersistentData{
		LastPrice:         lastPrice,
		LastChangePercent: lastChangePercent,
		LastUpdateTime:    lastUpdateTime,
	}

	file, err := os.Create(getDataFilePath())
	if err != nil {
		log.Printf("Error creating data file: %v", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		log.Printf("Error encoding data: %v", err)
	}
}

func loadPersistentData() PersistentData {
	var data PersistentData

	file, err := os.Open(getDataFilePath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error opening data file: %v", err)
		}
		return data
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error decoding data: %v", err)
	}

	return data
}

func getDataFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home directory: %v", err)
		return DATA_FILE
	}
	return filepath.Join(homeDir, DATA_FILE)
}

func getLockFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home directory: %v", err)
		return LOCK_FILE
	}
	return filepath.Join(homeDir, LOCK_FILE)
}