FROM golang:1.21.6-bookworm AS builder

WORKDIR /app
COPY ./ ./
RUN go build -ldflags="-s -w" -trimpath -o /app/main

FROM debian:bookworm

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/main /app/main

CMD ["/app/main"]
