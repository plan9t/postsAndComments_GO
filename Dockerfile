FROM golang:1.22 as builder
LABEL authors="planet-9"
# Установи рабочую директорию в контейнере
WORKDIR /app

# Копируем файлы модуля Go и ставим зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Сборка приложение
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ozon_server .

# Используй образ alpine для финального образа из-за его малого размера
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированный исполняемый файл из предыдущего степа
COPY --from=builder /app/ozon_server .

# Открываем порт, который использует приложение
EXPOSE 8090

# Запуск приложения
CMD ["./ozon_server"]