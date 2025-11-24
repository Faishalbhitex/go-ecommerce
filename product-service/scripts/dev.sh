#!/usr/bin/env bash

set -e

GREEN="\e[32m"
YELLOW="\e[33m"
RESET="\e[0m"

echo -e "${GREEN}[+] Start PostgreSQL...${RESET}"
pg_ctl -D $HOME/pg-go start
sleep 2
echo ""

echo -e "${GREEN}[+] Start Go Server...${RESET}"
mkdir -p bin
go build -o bin/product-service ./cmd/main.go
./bin/product-service &
echo ""

echo -e "${YELLOW}PostgreSQL and Go running on Background${RESET}"
echo -e "${YELLOW}Stop Server PostgreSQL and Go: ./scripts/stop.sh${RESET}"
