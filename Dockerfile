# Используйте официальный образ Go
FROM golang:1.22.1-alpine

# Установите рабочую директорию в контейнере
WORKDIR /app

# Копируйте go.mod и go.sum
COPY go.mod go.sum ./

# Загрузите зависимости
RUN go mod download

# Копируйте остальные части проекта
COPY cmd ./cmd
COPY configs ./configs
COPY deployments ./deployments
COPY internal ./internal
# Копируйте пользовательский конфигурационный файл из хоста в контейнер

# Сборка приложения
RUN go build -o ./ ./cmd/app/main.go
# RUN go build -o ./deployments/db/ ./deployments/db/deploy.go

# Запуск приложения
CMD ["./main"]
