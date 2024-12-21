## Порядок запуска проекта
- Postgres + Golang APP (updater) [yanlex-wow-guild-postgres + yanlex-wow-guild-updater](https://github.com/Yanlex/wow-guild-db )
- Front + Nginx - [wow-guild-front-nginx](https://github.com/Yanlex/wow-guild-front-nginx)
- ExpressJS API [yanlex-wow-guild-api](https://github.com/Yanlex/wow-guild-api-js )

# Задача приложения
Создание БД и работа с БД, в том числе обновление данных об игроках гильдии через API Raider io


### Docker контейнеры
Создаем сеть в которой наши контейнеры будут общаться
`docker network create wowguild`

## Настраиваем конфигурацию приложения
Настройки БД находятся в postgres-docker.yml
    POSTGRES_USER: user-name
    POSTGRES_PASSWORD: strong-password

Натросйки приложения находятся в файле backend-docker.yml
Дефолтные переменные
    DB_NAME: kvd_guild
    GUILD_REGION: eu
    GUILD_REALM: howling-fjord
    GUILD_NAME: "Ключик в дурку"
    DB_USER: user-name
    DB_PASS: strong-password
    DB_NETWORK: wowguild
    HOST_DB_PORT: 5432

### Настройка интервала обновления данных
В файле main.go
`_, _ = s.Every(1).Day().At("16:42").Do(updatePlayersHandler)`

### Логирование в файлы
Используем os.UserHomeDir() и основной путь /kvd/logs/ т.е создаем папку kvd в домашнем котологе пользователя куда будут писаться логи

### Запуск БД
`docker compose -f postgres-docker.yml up -d`

### Запуск приложения
Собираем проект, запускает backend-docker.yml который в свою очередь билдит backend.Dockerfile
`docker compose -f backend-docker.yml up -d --build`
