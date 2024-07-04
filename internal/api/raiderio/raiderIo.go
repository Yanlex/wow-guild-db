package fetch

import (
	"io"
	config "kvd/configs"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func init() {
	config.InitConfigDB()
}

func FetchRaiderIo() string {

	// URL по кторому получаем данные
	// КВД https://raider.io/api/v1/guilds/profile?region=eu&realm=howling-fjord&name=%D0%9A%D0%BB%D1%8E%D1%87%D0%B8%D0%BA%20%D0%B2%20%D0%B4%D1%83%D1%80%D0%BA%D1%83&fields=members
	url := viper.GetString("guild.raiderio_api_url")

	// Гет запрос
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Читаем данные
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Преобразование в строку
	bodyStr := string(body)
	// fmt.Println(bodyStr)

	// // Создать файл
	// file, err := os.Create("fetch.json")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer file.Close()

	// // Копируем все данные в файл
	// _, err = io.Copy(file, resp.Body)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// Читаем данные из файла
	// body, err := os.ReadFile("fetch.json")

	// result := gjson.Get(bodyStr, "members.#(character.name==\"Коррозийный\").character")
	// sd := gjson.Get(bodyStr, "members.#(character.name==\"Коррозийный\").character.profile_url")
	// fmt.Println(sd.String())

	return bodyStr
	// Блокировка выполнения программы, чтобы она не завершалась
}
