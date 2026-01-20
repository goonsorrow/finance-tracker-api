# ==========================================
# 1. Stage Builder (Собираем приложение)
# ==========================================
 FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o finance-tracker cmd/main.go

# ==========================================
# 2. Stage Runner (Запускаем приложение)
# ==========================================
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add bash curl

COPY --from=builder /app/finance-tracker .

COPY --from=builder /app/configs ./configs

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

# Запускаем приложение
CMD ["./finance-tracker"]
