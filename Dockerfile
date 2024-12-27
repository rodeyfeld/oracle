FROM golang:1.23 AS builder
WORKDIR /oracle
COPY go.mod go.sum ./
RUN go mod download

COPY . . 

WORKDIR /oracle/cmd
RUN go build -o /oracle/app-binary


FROM debian:bookworm-slim AS runner
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /oracle
COPY --from=builder /oracle/app-binary .
ENTRYPOINT ["/oracle/app-binary"]
