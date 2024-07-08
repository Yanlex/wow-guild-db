# Задача приложения
Создание БД и работе с БД, в том числе обновление данных об игроках гильдии через API Raider io

### Настройка интервала обновления данных
В файле main.go
`_, _ = s.Every(1).Day().At("16:42").Do(updatePlayersHandler)`

### Логирование в файлы
Используем os.UserHomeDir() и основной путь /kvd/logs/ т.е создаем папку kvd в домашнем котологе пользователя куда будут писаться логи

### Docker контейнеры
Создаем сеть в которой наши контейнеры будут общаться
`docker network create wowguild`

Запуск БД
`docker run --name yanlex-wow-guild-postgres --network wowguild -e POSTGRES_USER=user-name -e POSTGRES_PASSWORD=strong-password -v yanlex-wow-guild-postgres:/var/lib/postgresql/data -p 5432:5432 -d postgres:latest`

Запуск приложения
`docker build -t yanlex-wow-guild-updater .`
`docker run --network wowguild -d --name yanlex-wow-guild-updater -v yanlex-wow-guild-db-updater:/var/lib/postgresql/data yanlex-wow-guild-updater`

### Настраиваем конфигурацию приложения
:warning: Проверить настройки в configs/db.yaml
- в переменной `raiderio_api_url` должна быть ссылка на вашу гильдию
- в переменной `url` там где `"postgres://user-name:strong-password@localhost:5432"` нужно заменить `user-name:strong-password@localhost` на актуальные из команды запуска