# OABTray - Bitcoin Price Tracker 🚀🚀🚀🚀

OABTray is a simple and fun Bitcoin price tracker that sits in your system tray. It provides real-time updates on Bitcoin's price and percentage change, with amusing animations for significant price movements.

## Features

- Real-time Bitcoin price updates from CoinGecko API
- Customizable ALABA_FACTOR to set thresholds for price change reactions
- Fun animations for significant price movements
- Persistent data storage between sessions
- Easy-to-use system tray interface

## Installation

### Building from source

1. Ensure you have Go installed on your system.
2. Clone this repository:

   ```bash
   git clone https://github.com/ezeql/oabtray.git
   ```

3. Navigate to the project directory:

   ```bash
   cd oabtray
   ```

4. Build the application:

   ```bash
   go build
   ```

5. Run the executable:

   ```bash
   ./oabtray
   ```

### Brew

Brew formula is on /Users/ezeql/dev/homebrew-personal/Formula

## Usage

Once running, OABTray will appear in your system tray with the Bitcoin symbol (₿). The tray icon will display the current Bitcoin price and 24-hour percentage change.

- Click on the tray icon to see more options.
- You can adjust the ALABA_FACTOR, which determines the threshold for special animations.
- The app updates the price every 5 minutes.

## Configuration

The ALABA_FACTOR can be set to one of the following values: 0.5, 1, 1.5, 2, 2.5, or 3. This factor determines how sensitive the app is to price changes:

- When the price change percentage exceeds the positive ALABA_FACTOR, you'll see a "ALABADOOOOOOOO" animation.
- When it falls below the negative ALABA_FACTOR, you'll see a "PUTA MADRE" animation.

## Dependencies

This project uses the following external libraries:

- github.com/getlantern/systray

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

## Disclaimer

This application is for entertainment purposes only. Always do your own research before making any investment decisions.
