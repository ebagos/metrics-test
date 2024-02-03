FROM golang:1.21.6-bookworm AS builder

WORKDIR /app
COPY ./ ./
RUN go build -ldflags="-s -w" -trimpath -o /app/main

FROM debian:bookworm

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main"]
