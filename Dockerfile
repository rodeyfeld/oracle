FROM golang:1.23-alpine AS builder

EXPOSE 7777
WORKDIR /app

RUN apk update && apk add bash

COPY . /app 

CMD ["go", "run", "main.go"]
