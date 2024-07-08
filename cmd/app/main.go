package main

import (
	"fmt"
	config "kvd/configs"
	deploy "kvd/deployments/db"
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
	fmt.Println("DB UPDATER STARTED")

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
	_, _ = s.Every(1).Day().At("02:30").Do(updatePlayersHandler)

	// Запуск планировщика асинхронно
	s.StartAsync()
}

// var err error

func main() {

	timerDeploy := make(chan bool)
	timerMplus := make(chan bool)

	go func() {
		time.Sleep(10 * time.Second)
		timerDeploy <- true
	}()

	// Создаем канал для сигналов
	signals := make(chan os.Signal, 1)
	// Регистрируем канал для получения сигналов
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-timerDeploy
	deploy.Deploy()

	go func() {
		time.Sleep(25 * time.Second)
		timerMplus <- true
	}()

	<-timerMplus
	update.UpdateAllPlayers()
	// Блокируемся до получения сигнала
	sig := <-signals
	fmt.Println("Received signal:", sig)
}
