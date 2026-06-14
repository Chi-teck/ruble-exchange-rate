````# Ruble Exchange Rate

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
rer 100 EUR

# Rate for a specific date (DD.MM.YYYY).
rer --date 30.04.2020

# Inverse: convert 15000 rubles into US dollars.
rer --invert 15000 USD

# Raw numeric output (no decoration), handy for scripts.
rer --raw 100 EUR
```

Arguments are `[amount] [currency]`, defaulting to `1 USD`. The currency code is
case-insensitive.

## History

- `1.x` — original C implementation (branch `1.x`).
- `2.x` — current Go rewrite (branch `2.x`).

## License
GNU General Public License, version 2 or later.
````
