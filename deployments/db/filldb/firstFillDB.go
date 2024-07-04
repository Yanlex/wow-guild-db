package filldb

// #### Этот код запускается из db_createStructure.go
// #### Запускется один раз после создания базы данных для заполнения базы данных стартовой информацией о гильдии с АПИ Raider.IO

import (
	"context"
	"fmt"
	config "kvd/configs"
	fetch "kvd/internal/api/raiderio"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var pool *pgxpool.Pool
var err error

func init() {
	config.InitConfigDB()
	var ctx = context.Background()
	dbUrl := viper.GetString("db.urlKvd")
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}
	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
}

type Player struct {
	rank              string
	name              string
	guild             string
	realm             string
	race              string
	class             string
	gender            string
	faction           string
	achievementPoints string
	profileURL        string
	profileBanner     string
}

func FirstFillDB() {
	// Получение пути к домашнему каталогу
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	logFilePath := fmt.Sprintf("%s/kvd/logs/deploy/deploy.log", homeDir)

	// Создание всех необходимых каталогов, если они еще не существуют
	err = os.MkdirAll(fmt.Sprintf("%s/kvd/logs/deploy", homeDir), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем логирование в файл logs/update/updatePlayers.log
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := log.New(file, "[FirstFillDB] ", log.LstdFlags|log.Lshortfile)

	ctx := context.Background()
	resp := fetch.FetchRaiderIo()
	if resp == "" {
		logger.Println("Failed to fetch data from API")
		log.Fatalf("Failed to fetch data from API")
	}
	if err != nil {
		logger.Printf("Unable to create pool: %v", err)
		log.Fatalf("Unable to create pool: %v", err)
	}
	defer pool.Close()

	// Имя гильдии
	name := gjson.Get(resp, "name").String()
	if name == "" {
		logger.Println("Failed to extract name from API response")
		log.Fatalf("Failed to extract name from API response")
	}

	// Фракция
	faction := gjson.Get(resp, "faction").String()
	// if faction == "" {
	// 	logger.Println("Failed to extract faction from API response")
	// 	log.Fatalf("Failed to extract faction from API response")
	// }

	// Регион
	region := gjson.Get(resp, "region").String()
	// if region == "" {
	// 	logger.Println("Failed to extract region from API response")
	// 	log.Fatalf("Failed to extract region from API response")
	// }

	// Реалм
	realm := gjson.Get(resp, "realm").String()
	// if realm == "" {
	// 	logger.Println("Failed to extract realm from API response")
	// 	log.Fatalf("Failed to extract realm from API response")
	// }

	// Профиль
	profile_url := gjson.Get(resp, "profile_url").String()
	// if profile_url == "" {
	// 	logger.Println("Failed to extract profile_url from API response")
	// 	log.Fatalf("Failed to extract profile_url from API response")
	// }

	// Вставка данных в таблицу guild
	_, err = pool.Exec(ctx, `
        INSERT INTO guild (name, faction, region, realm, profile_url, created_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
    `, name, faction, region, realm, profile_url)
	if err != nil {
		logger.Printf("Failed to insert data: %v", err)
		log.Fatalf("Failed to insert data: %v", err)
	}

	defer fillPlayers(resp, file)
}

func fillPlayers(resp string, file *os.File) {
	logger := log.New(file, "[fillPlayers] ", log.LstdFlags|log.Lshortfile)

	totalMembers := gjson.Get(resp, "members.#")

	semaphoreBD := make(chan struct{}, 10)
	wg := sync.WaitGroup{}

	for i := 0; i < int(totalMembers.Int()); i++ {
		wg.Add(1)
		// Приведение i к int64
		// Создание пути с использованием fmt.Sprintf иначе gjson.Get выдаст ошибку too many arguments in call to gjson.Get
		rankPath := fmt.Sprintf("members.%d.rank", i) // Создание пути с использованием fmt.Sprintf
		rank := gjson.Get(resp, rankPath)

		namePath := fmt.Sprintf("members.%d.character.name", i)
		name := gjson.Get(resp, namePath)

		guild := "ключик в дурку"

		realmPath := fmt.Sprintf("members.%d.character.realm", i)
		realm := gjson.Get(resp, realmPath)

		racePath := fmt.Sprintf("members.%d.character.race", i)
		race := gjson.Get(resp, racePath)

		classPath := fmt.Sprintf("members.%d.character.class", i)
		class := gjson.Get(resp, classPath)

		genderPath := fmt.Sprintf("members.%d.character.gender", i)
		gender := gjson.Get(resp, genderPath)

		factionPath := fmt.Sprintf("members.%d.character.faction", i)
		faction := gjson.Get(resp, factionPath)

		achievementPointsPath := fmt.Sprintf("members.%d.character.achievement_points", i)
		achievement_points := gjson.Get(resp, achievementPointsPath)

		profileURLPath := fmt.Sprintf("members.%d.character.profile_url", i)
		profile_url := gjson.Get(resp, profileURLPath)

		profileBannerPath := fmt.Sprintf("members.%d.character.profile_banner", i)
		profile_banner := gjson.Get(resp, profileBannerPath)

		player := Player{
			rank:              rank.String(),
			name:              name.String(),
			guild:             guild,
			realm:             realm.String(),
			race:              race.String(),
			class:             class.String(),
			gender:            gender.String(),
			faction:           faction.String(),
			achievementPoints: achievement_points.String(),
			profileURL:        profile_url.String(),
			profileBanner:     profile_banner.String(),
		}
		// Логируем полученного игрока
		defer logger.Println(player)

		go func(p Player) {
			defer wg.Done()
			semaphoreBD <- struct{}{}
			insertObject(p)
			defer func() { <-semaphoreBD }()
		}(player)
	}
	defer fmt.Println("<->-- First fill DB DONE --<->")
	wg.Wait()
}

func insertObject(p Player) {
	ctx := context.Background()
	// Вставка данных в таблицу members
	_, err = pool.Exec(ctx, `
        INSERT INTO members (rank, name, guild, realm, race, class, gender, faction, achievement_points, profile_url, profile_banner, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP)
    `, p.rank, p.name, p.guild, p.realm, p.race, p.class, p.gender, p.faction, p.achievementPoints, p.profileURL, p.profileBanner)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

}
