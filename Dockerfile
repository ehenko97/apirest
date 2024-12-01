# Базовый образ
FROM ubuntu:latest

# Метаданные образа
LABEL authors="iGamez"

# Установка рабочего каталога
WORKDIR /Projectapirest

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода приложения
COPY . .

# Сборка приложения
RUN go build -o Projectapirest Projectapirest.go

# Указываем точку входа для запуска приложения
ENTRYPOINT ["./Projectapirest"]

# Указываем, что приложение слушает на порту 8080
EXPOSE 8080