package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/getlantern/systray"
)

const ALABA_FACTOR = 0.5
const INITIAL_DISPLAY_DURATION = 5 * time.Second
const SCREEN_WIDTH = 20
const ANIMATION_SPEED = 100 * time.Millisecond
const UPDATE_INTERVAL = 5 * time.Minute
const DATA_FILE = "bitcoin_tracker_data.gob"

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

var isFirstUpdate = true
var mu sync.Mutex
var lastPrice float64
var lastChangePercent float64
var lastUpdateTime time.Time

func main() {
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
	mBitcoinPrice := systray.AddMenuItem("Loading...", "Bitcoin price")
	mBitcoinPrice.Disable()
	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	if lastPrice > 0 {
		updateTray(lastPrice, lastChangePercent)
		log.Printf("Loaded price from disk: $%.2f (%+.2f%%)", lastPrice, lastChangePercent)
	}
}

func onExit() {
	savePersistentData()
}

func priceUpdater() {
	// Esperar un poco para asegurarse de que la interfaz de usuario esté lista
	time.Sleep(1 * time.Second)

	// Si los datos son antiguos o no existen, actualizar inmediatamente
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
	price, changePercent, err := fetchPriceFromCoinGecko()
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

func fetchPriceFromCoinGecko() (float64, float64, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&include_24hr_change=true")
	if err != nil {
		return 0, 0, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	var data CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("JSON decode error: %v", err)
	}

	return data.Bitcoin.USD, data.Bitcoin.USD24HChange, nil
}

func updateTrayWithInitialDisplay(price, changePercent float64) {
	priceStr := formatPriceString(price, changePercent)
	systray.SetTitle(priceStr)
	systray.SetTooltip(priceStr)

	time.Sleep(INITIAL_DISPLAY_DURATION)

	updateTray(price, changePercent)
}

func updateTray(price, changePercent float64) {
	priceStr := formatPriceString(price, changePercent)
	systray.SetTitle(priceStr)
	systray.SetTooltip(priceStr)

	if math.Abs(changePercent-lastChangePercent) >= ALABA_FACTOR {
		var animationText string
		if changePercent >= lastChangePercent+ALABA_FACTOR {
			animationText = strings.Repeat("ALABADOOOOOOOO ", 5)
		} else if changePercent <= lastChangePercent-ALABA_FACTOR {
			animationText = strings.Repeat("PUTA MADRE ", 5)
		}

		if animationText != "" {
			go animateTrainSign(animationText)
		}
	}
}

func formatPriceString(price, changePercent float64) string {
	var emoji string
	if changePercent >= ALABA_FACTOR {
		emoji = "🟢"
	} else if changePercent < -ALABA_FACTOR {
		emoji = "🔴"
	} else {
		emoji = "⚪"
	}

	rockets := getRocketEmojis(changePercent)
	return fmt.Sprintf("₿ %s $%.2f (%+.2f%%) %s", emoji, price, changePercent, rockets)
}

func getRocketEmojis(changePercent float64) string {
	rocketCount := int(math.Abs(changePercent)) % int(ALABA_FACTOR*10)
	return strings.Repeat("🚀", rocketCount)
}

func displayError(err error) {
	errorMsg := fmt.Sprintf("₿ Error: %s", err.Error())
	systray.SetTitle(errorMsg)
	systray.SetTooltip(errorMsg)
}

func animateTrainSign(text string) {
	textLength := len(text)

	for i := 0; i < textLength; i++ {
		displayText := text[i:] + text[:i]
		displayText = displayText[:SCREEN_WIDTH]
		systray.SetTitle(displayText)
		time.Sleep(ANIMATION_SPEED)
	}

	// Restaurar el precio después de la animación
	systray.SetTitle(formatPriceString(lastPrice, lastChangePercent))
}

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
