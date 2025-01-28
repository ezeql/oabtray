# 🚀🚀🚀 OAB Tray 🚀🚀🚀


A system tray application that tracks and displays the current Bitcoin price in USD, with additional features for price change notifications and animations.

![alt text](image-1.png)

## Features

- 📈 Real-time Bitcoin price tracking from Binance API
- 🖥️ System tray display with price and percentage change
- 🎬 Price change animations for significant movements (≥5%)
- 🚀 "To the Mow-n" mode to display price in millions EX: $0.100M

## Installation

1. Using Homebrew (macOS):

   ```bash
   # Install
   brew tap ezeql/personal 
   brew install oabtray

   # Run
   oabtray # Ctrl+C to exit
   ```

2. Optional: To run OAB Tray automatically when your system starts:

    ```bash
    brew services start oabtray
    ```

## Usage

- The application runs in the system tray
- Click the tray icon to see additional options:
  - Current price and version information
  - Toggle "To the Moon" mode
  - Quit the application
- Significant price changes (≥5%) trigger animations:
  - 🚀 for positive changes
  - 🧂 for negative changes

## Configuration

The application automatically saves its state to

- `~/bitcoin_tracker_data.gob` - Last known price and settings
- `~/bitcoin_tracker.lock` - Lock file to prevent multiple instances

## Dependencies

- github.com/getlantern/systray
- github.com/gofrs/flock

## License

MIT License
