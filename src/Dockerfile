FROM golang:1.15.2 AS builder

RUN mkdir -p /app
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:3.12.3

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app/app ./
RUN apk add --no-cache bash

ENTRYPOINT ["./app"]