package main

import (
	"fmt"
	config "kvd/configs"
	"kvd/internal/db/update"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
)

func updatePlayersHandler() {
	fmt.Println("CRON TASK STARTED")
	update.UpdateAllPlayers()
}

func init() {
	config.InitConfigDB()

	// Крон планировщик
	// Загрузка локации
	est, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		// Логирование ошибки вместо паники
		log.Printf("Error loading location: %v", err)
		return
	}

	// Создание планировщика с указанной локацией
	s := gocron.NewScheduler(est)

	// Планирование задачи
	_, _ = s.Every(1).Day().At("13:42").Do(updatePlayersHandler)

	// Запуск планировщика асинхронно
	s.StartAsync()
}

// var err error

func main() {
	fmt.Println("DB UPDATER STARTED")

	// Создаем канал для сигналов
	signals := make(chan os.Signal, 1)
	// Регистрируем канал для получения сигналов
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Блокируемся до получения сигнала
	sig := <-signals
	fmt.Println("Received signal:", sig)
}
