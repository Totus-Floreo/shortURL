#!/bin/bash

pg_url=postgres:password@localhost:32773/links
httpport=:3011
gRPCport=:3022

export pg_url httpport gRPCport

go run ./cmd/shortURL/main.go -dbType pgx