# Задача приложения
Создание БД и обновление данных об игроках гильдии через API Raider io по крон задаче раз в день.

### Настройка интервала обновления данных
В файле main.go
`_, _ = s.Every(1).Day().At("16:42").Do(updatePlayersHandler)`

### Логирование в файлы
Используем os.UserHomeDir() и основной путь /kvd/logs/ т.е создаем папку kvd в домашнем котологе пользователя куда будут писаться логи

### Docker контейнеры
Запускаем **docker-compose.yml** из папки deployments/docker/
>:warning: Важно указать сильный логин и пароль в полях 
>**POSTGRES_USER** и **POSTGRES_PASSWORD**
>**PGADMIN_DEFAULT_EMAIL** и **PGADMIN_DEFAULT_PASSWORD**
- запуск находясь в котологе с файлом
`docker compose up -d`

### Настраиваем конфигурацию приложения
:warning: Проверить настройки в configs/db.yaml
- в переменной `raiderio_api_url` должна быть ссылка на вашу гильдию
- в переменной `url` там где `"postgres://user-name:strong-password@localhost:5432"` нужно заменить user-name:strong-password на актуальные из docker-compose.yml

:warning: Файл configs/config.go 
- важно указать правильный путь до каталога `configs` в переменной `viper.AddConfigPath("$HOME/goproject/wow-guild-website/configs") `

### Запускаем скрипт создания структуры БД
Скрипт **deploy.go** лежит в папке deployments/db/
`go run deploy.go`

### Запуск приложения в менеджере процессов PM2
- глобально ставим pm2
`npm install pm2 -g`
-  после всех настроек собираем приложение
`go build main.go`
- запускаем приложение
`pm2 start ./main`