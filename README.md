# Ruble Exchange Rate

A simple command-line tool to obtain ruble exchange rates from the CBR (Central
Bank of the Russian Federation).

## Installation

### Binary

You can download the binary from the releases page on GitHub and add it to your $PATH.

### From source

Ensure that you have a supported version of Go properly installed and set up.

```bash
go install github.com/Chi-teck/ruble-exchange-rate/v2/cmd/rer@latest
```

## Usage

```bash
# Latest rate (defaults to 1 USD).
rer

# Convert a specific amount of a currency into rubles.
rer -a 100 -c EUR

# Rate for a specific date (DD.MM.YYYY).
rer -d 30.04.2020

# Inverse conversion of 15000 rubles into US dollars.
rer -i -a 15000 -c USD

# Raw numeric output (no decoration), handy for scripts.
rer -r -a 100 -c EUR
```

Options (each has a short alias). With no flags, `rer` defaults to `1 USD`:

- `-a, --amount` — amount to convert (default `1`).
- `-c, --currency` — currency code, e.g. `USD` or `EUR` (default `USD`, case-insensitive).
- `-d, --date` — rate for a specific date, `DD.MM.YYYY` (default: latest).
- `-i, --invert` — show the inverse (convert rubles into the currency).
- `-r, --raw` — print only the number, no decoration.

## History

- `1.x` — original C implementation.
- `2.x` — current Go rewrite.

## License

GNU General Public License, version 2 or later.
