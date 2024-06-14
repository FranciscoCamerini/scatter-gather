#!/bin/bash

ORCHESTRATOR_PORT=8080
WORKER_PORTS=(8081 8082 8083 8084)

BASE_DIR=$(dirname "$0")
STORAGE_DIR="$BASE_DIR/storage" # Used for storing PID's
PUBLIC_DIR="$BASE_DIR/public" # Contains all text files for parsing

mkdir -p "$STORAGE_DIR"

make -C "$BASE_DIR"

for port in "${WORKER_PORTS[@]}"; do
    "$BASE_DIR/worker/worker" -port="$port" -publicdir="$PUBLIC_DIR" -pidfile="$STORAGE_DIR/.worker-$port.pid" & # Spawn workers in the background
done

"$BASE_DIR/orchestrator/orchestrator" -port="$ORCHESTRATOR_PORT" -pidfile="$STORAGE_DIR/.orchestrator.pid" -workers=$(IFS=,; echo "${WORKER_PORTS[*]}")
