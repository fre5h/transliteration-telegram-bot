FROM golang:1.20.1-alpine AS builder
RUN mkdir /build
RUN mkdir /build/internal
ADD go.mod go.sum cmd/main.go /build/
ADD internal/ /build/internal/
WORKDIR /build
RUN go build -o transliteration-telegram-bot

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/transliteration-telegram-bot /app/
WORKDIR /app
CMD ["./transliteration-telegram-bot"]
