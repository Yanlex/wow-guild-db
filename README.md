## Порядок запуска проекта
- Postgres + Golang APP (updater) [yanlex-wow-guild-postgres + yanlex-wow-guild-updater](https://github.com/Yanlex/wow-guild-db )
- Front + Nginx - [wow-guild-front-nginx](https://github.com/Yanlex/wow-guild-front-nginx)

# Задача приложения
Создание БД и работа с БД, в том числе обновление данных об игроках гильдии через API Raider io
Собственное API к своей БД

### Docker контейнеры
Если сеть еще не создана, создаем
`docker network create wowguild`

## Настраиваем конфигурацию приложения
Настройки БД находятся в postgres-docker.yml
- POSTGRES_USER: user-name
- POSTGRES_PASSWORD: strong-password

Натросйки приложения находятся в файле backend-docker.yml
Дефолтные переменные
- DB_NAME: kvd_guild
- GUILD_REGION: eu
- GUILD_REALM: howling-fjord
- GUILD_NAME: "Ключик в дурку"
- DB_USER: user-name
- DB_PASS: strong-password
- DB_NETWORK: wowguild
- DB_ADDRESS: yanlex-wow-guild-postgres
- HOST_DB_PORT: 5432

### Настройка интервала обновления данных
Делает запросы к стороннему АПИ и сверяет есть ли изменения в данных игрока, например рейтинг м+ или количетсво ачивок.
В файле main.go
`_, _ = s.Every(1).Day().At("19:12").Do(updatePlayersHandler)`

### Логирование в файлы
Это скорее просто для опыта реализовано минимально.
Используем os.UserHomeDir() и основной путь /kvd/logs/ т.е создаем папку kvd в домашнем котологе пользователя куда будут писаться логи

### Запуск БД
`docker compose -f postgres-docker.yml up -d`

### Запуск приложения
Собираем проект, запускает backend-docker.yml который в свою очередь билдит backend.Dockerfile
`docker compose -f backend-docker.yml up -d --build`


## API
Работает на 3000 порту
Возвращает список игроков гильдии
/api/get-members
Возвращает Rank, Name, Mythic Rating, Guild, Class
/api/guild-data
Шарим папку с аватарками
/api/avatar/
Шарим папку с классами
/api/class/