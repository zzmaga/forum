# Используем официальный образ Go
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты для SQLite
RUN apk --no-cache add ca-certificates sqlite

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/main .

# Копируем HTML шаблоны и статические файлы
COPY --from=builder /app/ui ./ui

# Создаем директорию для базы данных
RUN mkdir -p /root/data

# Меняем владельца файлов
RUN chown -R appuser:appgroup /root/

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Команда для запуска
CMD ["./main"]
