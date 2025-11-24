#!/usr/bin/env bash

set -e

GREEN="\e[32m"
RESET="\e[0m"

echo -e "${GREEN}[+]Stopping Server Go..${RESET}"
pkill -f product-service || true
rm -f bin/product-service
echo ""

echo -e "${GREEN}[+]Stopping Server PostgreSQL..${RESET}"
pg_ctl -D $HOME/pg-go stop
echo ""

echo -e "${GREEN}[+]Go and PostgreSQL shutdown successfuly!${RESET}"
