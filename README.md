# Ruble Exchange Rate

A simple command line tool to obtain ruble exchange rates from CBR (The Central
Bank of the Russian Federation).

## Installation

### Prebuilt binaries

Download an archive for your platform from the
[Releases page](https://github.com/Chi-teck/ruble-exchange-rate/releases),
extract the `rer` binary, and place it somewhere on your `$PATH`:

```
tar -xzf rer_*_linux_amd64.tar.gz
sudo install rer /usr/local/bin/
```

Builds are available for Linux, macOS, and Windows (amd64 and arm64); the
Windows archive is a `.zip`.

### From source
```
git clone https://github.com/Chi-teck/ruble-exchange-rate.git
cd ruble-exchange-rate
go build ./cmd/rer   # produces ./rer
```
Or install it directly into your `$GOBIN`:
```
go install github.com/Chi-teck/ruble-exchange-rate/cmd/rer@latest
```

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
