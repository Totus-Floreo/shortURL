#!/bin/bash

pg_url=postgres:password@localhost:32768/links
PORT=3010
gRPCport=3011

export pg_url PORT gRPCport

go run ./cmd/shortURL/main.go -dbType pgx