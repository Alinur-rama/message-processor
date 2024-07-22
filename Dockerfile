FROM golang:1.22.3-alpine AS builder

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git

# Копирование и загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Финальный этап
FROM alpine:latest  

RUN apk --no-cache add ca-certificates bash curl

WORKDIR /root/

# Копирование исполняемого файла из этапа сборки
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запуск приложения
CMD ["/root/main"]