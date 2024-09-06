FROM golang:1.23-alpine AS builder

RUN apk update && apk add bash
EXPOSE 7777

RUN mkdir /app
COPY . /app

WORKDIR /app


CMD ["go", "run", "."]
