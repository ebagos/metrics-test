FROM golang:1.21.6-alpine3.19 AS builder

WORKDIR /app
COPY ./ ./
RUN go build -ldflags="-s -w" -trimpath -o /app/main

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main"]
