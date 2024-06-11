#!/bin/bash

SERVER_PORT=8080
CRAWLER_PORTS=(8081 8082 8083 8084)

make

./master/master -pidfile=storage/.master.pid -crawlers=$(IFS=,; echo "${CRAWLER_PORTS[*]}")&

for port in "${CRAWLER_PORTS[@]}"; do
    ./crawler/crawler -port=$port -pidfile=storage/.crawler-$port.pid&
done
