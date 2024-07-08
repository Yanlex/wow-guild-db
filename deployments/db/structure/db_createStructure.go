package structure

// #### С Создаем структуру БД, после запускаем заполнение стартовой информацией о гильдии с АПИ raid.io через defer firstfilldb.FirstFillDB()

import (
	"context"
	"fmt"
	config "kvd/configs"
	filldb "kvd/deployments/db/filldb"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func Init() {
	config.InitConfigDB()

	// Конфигурация подключения
	connConfig, err := pgx.ParseConfig(viper.GetString("db.url"))
	if err != nil {
		log.Fatalf("Configuration parsing error: %v\n", err)
	}
	dbBuild(connConfig)
}

func dbBuild(connConfig *pgx.ConnConfig) {
	// Получение пути к домашнему каталогу
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	logFilePath := fmt.Sprintf("%s/kvd/logs/deploy.log", homeDir)

	// Создание всех необходимых каталогов, если они еще не существуют
	err = os.MkdirAll(fmt.Sprintf("%s/kvd/logs", homeDir), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем логирование в файл logs/update/updatePlayers.log
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "[DEPLOY] ", log.LstdFlags|log.Lshortfile)

	// Подключение
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v\n", err)
		logger.Fatalf("Error connecting to PostgreSQL: %v\n", err)
	}
	// Закрыть соединение после выполнения функции
	defer conn.Close(context.Background())

	dbName := viper.GetString("dbCreateStructure.dbName")

	// Проверяем, существует ли база данных
	checkDBExistsQuery := fmt.Sprintf("SELECT datname FROM pg_database WHERE datname = '%s'", dbName)
	var existsDB string // Изменено с bool на string
	err = conn.QueryRow(context.Background(), checkDBExistsQuery).Scan(&existsDB)
	if err != nil && err != pgx.ErrNoRows {
		log.Fatalf("Error checking if database exists: %v\n", err)
		logger.Fatalf("Error checking if database exists: %v\n", err)
	}

	if existsDB == "" { // Проверяем, что exists пустая строка, что означает отсутствие базы данных
		// SQL запрос для создания БД
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)

		// Выполняем запрос
		_, err = conn.Exec(context.Background(), createDBQuery)
		if err != nil {
			log.Fatalf("Failed to create database: %v\n", err)
			logger.Fatalf("Failed to create database: %s %v\n", dbName, err)
		}
		logger.Printf("The database: %s has been successfully created", dbName)
		log.Printf("The database: %s has been successfully created\n", dbName)
	} else {
		logger.Printf("The database: %s already exists", dbName)
		log.Printf("The database: %s already exists\n", dbName)
		// Удаление комментариев о необходимости удаления базы данных перед созданием, так как это не требуется
	}

	// Подключение к базе данных kvd_guild
	connConfig.Database = viper.GetString("dbCreateStructure.dbName")
	conn, err = pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		logger.Fatalf("Error connecting to a new database: %v\n", err)
		log.Fatalf("Error connecting to a new database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// SQL запрос для создания таблиц и столбцов в базе kvd_guild
	createTableAndRow := `
	CREATE TABLE IF NOT EXISTS guild (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		faction VARCHAR(255) ,
		region VARCHAR(255) ,
		realm VARCHAR(255),
		profile_url VARCHAR(255) ,
		created_at TIMESTAMP DEFAULT now()
	);
	CREATE TABLE IF NOT EXISTS members (
    id SERIAL PRIMARY KEY,
    rank INTEGER,
    name VARCHAR(255) NOT NULL,
    mythic_plus_scores_by_season INTEGER DEFAULT 0,
    guild VARCHAR(255),
    realm VARCHAR(255) DEFAULT '',
    race VARCHAR(255),
    class VARCHAR(255),
    gender VARCHAR(255),
    faction VARCHAR(255),
    achievement_points INTEGER,
    profile_url VARCHAR(255),
    thumbnail_url VARCHAR(255) DEFAULT '',
    profile_banner VARCHAR(255),
    created_at TIMESTAMP DEFAULT now()
);

	`

	_, err = conn.Exec(context.Background(), createTableAndRow)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
		logger.Fatalf("Failed to create table: %v\n", err)
	}
	defer filldb.FirstFillDB()
}
