# OABTray - Bitcoin Price Tracker 🚀🚀🚀🚀

[![Go Multi-Platform Build and Release](https://github.com/ezeql/oabtray/actions/workflows/go.yml/badge.svg)](https://github.com/ezeql/oabtray/actions/workflows/go.yml)

OABTray is a simple and fun Bitcoin price tracker that sits in your system tray. It provides real-time updates on Bitcoin's price and percentage change from Binance, with amusing animations for significant price movements.

## Features

- Real-time Bitcoin price updates from Binance API
- Price display in USD with optional millions mode
- Dynamic emoji indicators: 🟢 (up), 🔴 (down), ⚪ (no change)
- Rocket 🚀 and salt 🧂 indicators based on price movements
- Fun animations for significant price changes (±5%)
- Persistent data storage between sessions
- Easy-to-use system tray interface
- 30-second price update interval

## System Requirements

- macOS (currently macOS-only)
- Go 1.23.3 or later (for building from source)

## Installation

### Building from Source

1. Ensure you have Go 1.23.3 or later installed on your macOS system
1. Clone this repository:

```bash
git clone https://github.com/ezeql/oabtray.git
```

1. Navigate to the project directory:

```bash
cd oabtray
```

1. Build the application:

```bash
./build.sh
```

1. Run the executable:

```bash
./oabtray
```

### Brew

Brew formula is on /Users/ezeql/dev/homebrew-personal/Formula

## Usage

Once running, OABTray will appear in your system tray with the Bitcoin symbol (₿). The tray icon will display:

- Current Bitcoin price (with thousands separator)
- 24-hour percentage change
- Emoji indicator (🟢, 🔴, or ⚪)
- Rocket 🚀 or salt 🧂 indicators based on price movement

### Options

- Click on the tray icon to see options
- Toggle "Set price in millions" to switch between normal and millions display mode
- View the current version
- Price updates every 30 seconds
- Special animations trigger when price change exceeds ±5%:
  - "ALABADO!!!" for +5% or more
  - "PUTA MADRE!" for -5% or less

## Price Indicators

The app shows price movement intensity with emojis:

- Up movements: Rocket emojis 🚀 (more rockets = bigger movement)
- Down movements: Salt emojis 🧂 (more salt = bigger drop)
- Current trend: 🟢 (up), 🔴 (down), or ⚪ (no change)

## Dependencies

This project uses the following external libraries:

- github.com/getlantern/systray
- github.com/gofrs/flock

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

## Disclaimer

This application is for entertainment purposes only. Always do your own research before making any investment decisions.
