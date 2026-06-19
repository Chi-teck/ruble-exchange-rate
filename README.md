# Ruble Exchange Rate

A simple command line tool to obtain ruble exchange rates from CBR (The Central
Bank of the Russian Federation).

## Installation
```
git clone https://github.com/Chi-teck/ruble-exchange-rate.git
cd ruble-exchange-rate
go build      # produces ./rer
```
Optionally run `go install` to place the `rer` binary in your `$GOBIN`.

## Usage
```
# Latest rate (defaults to 1 USD).
rer

# Convert a specific amount of a currency into rubles.
rer -a 100 -c EUR

# Rate for a specific date (DD.MM.YYYY).
rer --date 30.04.2020

# Inverse: convert 15000 rubles into US dollars.
rer --invert -a 15000 -c USD

# Raw numeric output (no decoration), handy for scripts.
rer --raw -a 100 -c EUR
```

Options (each has a short alias), defaulting to `1 USD`:

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
