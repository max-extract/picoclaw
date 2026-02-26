#!/usr/bin/env sh
set -eu

if [ "$#" -lt 2 ]; then
  echo "usage: wait-for-text.sh <file> <text>"
  exit 2
fi

FILE="$1"
TEXT="$2"
TIMEOUT="${3:-30}"

i=0
while [ "$i" -lt "$TIMEOUT" ]; do
  if [ -f "$FILE" ] && grep -q "$TEXT" "$FILE"; then
    echo "found"
    exit 0
  fi
  sleep 1
  i=$((i+1))
done

echo "timeout"
exit 1
