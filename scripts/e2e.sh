#!/usr/bin/env bash
#
# e2e.sh — smoke-test the `rer` binary against live CBR data.
#
# Builds ./rer, then runs a sequence of labelled scenarios (happy paths and
# error cases), printing the command, its output, and exit code for each so a
# human can eyeball the results. This is a manual aid, not an assertion-based
# test — see cmd/rer/main_test.go for the automated unit tests.
#
# Usage:
#   scripts/e2e.sh             # build, then run every scenario
#   scripts/e2e.sh --no-build  # reuse the existing ./rer
#
# Note: this hits the live feed at cbr.ru. cbr.ru is appended to NO_PROXY so
# direct egress works behind a filtering proxy (e.g. Claude's sandbox); it is
# harmless when no proxy is configured.

set -uo pipefail

cd "$(dirname "$0")/.."

BIN=./rer
export NO_PROXY="${NO_PROXY:+$NO_PROXY,}cbr.ru"

# Colours, only when stdout is a terminal and NO_COLOR is unset.
if [[ -t 1 && -z "${NO_COLOR:-}" ]]; then
  bold=$(tput bold); dim=$(tput dim); reset=$(tput sgr0)
else
  bold=; dim=; reset=
fi

if [[ "${1:-}" != "--no-build" ]]; then
  echo "${bold}Building $BIN ...${reset}"
  go build -o "$BIN" ./cmd/rer || { echo "build failed" >&2; exit 1; }
fi
[[ -x "$BIN" ]] || { echo "$BIN not found; run without --no-build" >&2; exit 1; }

n=0
# run "<what to expect>" [args...] — print the command, its output and exit code.
run() {
  local expect=$1; shift
  n=$((n + 1))
  printf '\n%s[%02d] $ %s %s%s\n' "$bold" "$n" "$BIN" "$*" "$reset"
  printf '%s     expect: %s%s\n' "$dim" "$expect" "$reset"
  "$BIN" "$@"
  printf '%s     [exit %d]%s\n' "$dim" "$?" "$reset"
}

echo
echo "${bold}=== Happy paths ===${reset}"
run "1 USD = <rate> RUB (latest date)"
run "100 EUR = <rate> RUB"                       -a 100 -c EUR
run "lowercase 'eur' accepted (1 EUR = <rate>)"  -c eur
run "inverse: 1 RUB = <small> EUR"               -i -c EUR
run "raw: bare number, 4 decimals, no text"      -r -c EUR
run "historical rate for 01.06.2024"             -c USD -d 01.06.2024

echo
echo "${bold}=== Error / edge cases (expect exit 1, friendly message) ===${reset}"
run "unknown currency, then lists available codes" -c XXX
run "amount must be positive"                      -a -5 -c USD
run "zero amount rejected"                         -a 0 -c USD
run "non-numeric amount rejected"                  -a abc -c USD
run "bad date format rejected"                     -c USD -d 2024-01-01
run "no data before 01.07.1992"                    -c USD -d 01.01.1990

echo
echo "${bold}=== Usage ===${reset}"
run "help text listing all five flags (exit 0)"    --help

echo
echo "${bold}Done: $n scenarios. Eyeball the output above.${reset}"
