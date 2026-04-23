# ЭТАП 1: СБОРКА
# Используем нужную версию Go на Alpine
FROM golang:1.26.2-alpine AS builder

WORKDIR /

# Копируем файлы с зависимостями
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /app/api ./cmd/api

# ЭТАП 2: ФИНАЛЬНЫЙ ОБРАЗ
# Используем чистый Alpine (или даже scratch)
FROM alpine:3.20

# Создаём НЕ-root пользователя для безопасности [citation:7]
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Копируем бинарник из образа-сборщика
COPY --from=builder /app/api /tasks-api

# Переключаемся на непривилегированного пользователя
USER appuser

CMD ["/tasks-api"]