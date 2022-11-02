# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main src/infrastructure/cmd/server/server.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/main /main

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/main"]