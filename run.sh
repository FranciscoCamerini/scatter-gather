#!/bin/bash

SERVER_PORT=8080
CRAWLER_PORTS=(8081 8082 8083 8084)

BASE_DIR=$(dirname "$0")
STORAGE_DIR="$BASE_DIR/storage" # Used for storing PID's

mkdir -p "$STORAGE_DIR"

make -C "$BASE_DIR"

"$BASE_DIR/master/master" -pidfile="$STORAGE_DIR/.master.pid" -crawlers=$(IFS=,; echo "${CRAWLER_PORTS[*]}") &

for port in "${CRAWLER_PORTS[@]}"; do
    "$BASE_DIR/crawler/crawler" -port="$port" -pidfile="$STORAGE_DIR/.crawler-$port.pid" &
done
