# Stage 1: Builder
FROM golang:1.23.2 AS builder

# Рабочая директория
WORKDIR /github.com/mskovv/tg-bot-subaru96

# Копируем зависимости и устанавливаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Установка Air
RUN go install github.com/air-verse/air@latest && \
    echo "Air установлен: $(which air)" && \
    ls -l /go/bin/air

# Stage 2: Final image
FROM golang:1.23.2

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y bash && apt-get clean

# Рабочая директория
WORKDIR /usr/src/app

# Копируем исходный код и Air
COPY --from=builder /github.com/mskovv/tg-bot-subaru96 .
COPY --from=builder /go/bin/air /usr/local/bin/air

# Проверяем наличие Air
RUN which air && chmod +x /usr/local/bin/air

# Открываем порт
EXPOSE 9000