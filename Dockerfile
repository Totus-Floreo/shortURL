# syntax=docker/dockerfile:1

FROM golang:1.19

ENV httpport=:8080
ENV gRPCport=:50051

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080

EXPOSE 50051

WORKDIR /app/cmd/shortURL

RUN CGO_ENABLED=0 GOOS=linux go build -o shorturl .

ENTRYPOINT ["./shorturl"]
