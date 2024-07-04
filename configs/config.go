package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfigDB() {
	viper.SetConfigName("db")                                        // Имя файла конфигурации без расширения
	viper.SetConfigType("yaml")                                      // Тип файла конфигурации
	viper.AddConfigPath("$HOME/goproject/wow-guild-website/configs") // Путь к директории с конфигурацией
	// Чтение конфигурации
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
	}
}
