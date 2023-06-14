# shorturl
Test task for Ozon Internship.

## Intern Developer Test Task
Task:
Implement a service that provides an API for creating shortened links.

The link should be:

    - Unique; only one shortened link should point to a single original URL.
    - 10 characters long.
    - Consist of lowercase and uppercase Latin alphabet characters, digits, and the underscore symbol.

The service should be written in Go and accept the following HTTP requests:

    - POST method, which will save the original URL in the database and return the shortened URL.
    - GET method, which will accept the shortened URL and return the original URL.

Extra Credit:
Make the service work via gRPC, i.e., create a proto file and implement the service with two corresponding endpoints.

The solution should meet the following conditions:

    - The service is distributed as a Docker image.
    - In-memory storage and PostgreSQL are expected as storage solutions. The storage to be used is specified as a parameter when starting the service.
    - The implemented functionality is covered by unit tests.

Result as a public repository on github.com.
###### Thanks ChatGPT
## Quickstart
### Params
```sh
#inmemory - cache db based on map
#pgx - postgresql db(scripts for db in folder script/sql/)
-dbType=<Type> #Optional, default use pgx
```
### Just Code, No More
Setting and run this script
```sh
#!/bin/bash

pg_url=postgres:password@localhost:32773/links # postgres url to connect
httpport=:3011 # http listener port(based on Gin)
gRPCport=:3022 # gRPC listener port 

export pg_url httpport gRPCport 

go run ./cmd/shortURL/main.go -dbType pgx
```
like this
```sh
./scripts/run.sh
```
### Docker Container and External DB
How to build
```sh
docker build . --tag shorturl
```
How to run
```sh
docker run -e pg_url=<pgx_user>:<pgx_password>@<pgx_host>:<pgx_port>/links /
-p <http_port>:8080 /
-p <gPRC_port>:50051 /
shorturl -dbType=<typeofdb> #flag is optional, default use pgx
```
Example 
```sh
docker run -e pg_url=postgres:password@192.168.1.2:32771/links -p 3010:8080 -p 3020:50051 shorturl -dbType=pgx
```
### Docker Compose with Postgres 
```sh
#flag -d for silent start
docker compose up -d
```
```sh
#port forwarding is not displayed in docker desktop
docker compose ps
```