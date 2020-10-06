# Ruble Exchange Rate

A simple command line tool to obtain ruble exchange rate from CBR (The Central Bank of the Russian Federation).

## Installation
```
git clone https://github.com/Chi-teck/rer.git
cd rer
make
make install
```

## Usage
```
# Get rate against multipe currencies.
./rer

# Get rate for a specific date.
./rer --date=30.04.2020

# Get rate against a specific currency.
./rer -c usd

# Convert 15000 rubbles into US dollars.
./rer -i -a 15000 -c usd
```

## License
GNU General Public License, version 2 or later.
