package update

import (
	"context"
	"fmt"
	"io"
	config "kvd/configs"
	fetch "kvd/internal/api/raiderio"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var pool *pgxpool.Pool
var ctx context.Context
var logger *log.Logger
var file *os.File

func init() {
	// Получаем конфигурацию соединения с БД
	config.InitConfigDB()
	dbUrl := viper.GetString("db.urlKvd")
	ctx = context.Background()
	// fmt.Println(dbUrl)
	connConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Configuration parsing error: %v\n", err)
	}
	// Создаем пул соединений
	pool, err = pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Получение пути к домашнему каталогу
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	logFilePath := fmt.Sprintf("%s/kvd/logs/update/updatePlayers.log", homeDir)

	// Создание всех необходимых каталогов, если они еще не существуют
	err = os.MkdirAll(fmt.Sprintf("%s/kvd/logs/update", homeDir), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем логирование в файл logs/update/updatePlayers.log
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "[UPDATEPlAYERS] ", log.LstdFlags|log.Lshortfile)
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
type PlayerDB struct {
	rank                         string
	name                         string
	mythic_plus_scores_by_season string
	guild                        string
	realm                        string
	race                         string
	class                        string
	gender                       string
	faction                      string
	achievementPoints            string
	profileURL                   string
	profileBanner                string
}

func UpdateAllPlayers() {
	fmt.Println("UPDATE PLAYERS STARTED")

	// Получаем данные из API
	resp := fetch.FetchRaiderIo()
	if resp == "" {
		log.Fatalf("Failed to fetch data from API")
	}

	// Получаем из Базы данных таблицу members
	rows, err := pool.Query(context.Background(), "SELECT  rank, name, mythic_plus_scores_by_season,  guild, realm, race, class, gender, faction, achievement_points, profile_url, profile_banner FROM members")
	if err != nil {
		log.Fatalf("Query error: %v\n", err)
	}
	defer rows.Close()

	var players []PlayerDB
	for rows.Next() {

		var player PlayerDB
		if err := rows.Scan(&player.rank, &player.name, &player.mythic_plus_scores_by_season, &player.guild, &player.realm, &player.race, &player.class, &player.gender, &player.faction, &player.achievementPoints, &player.profileURL, &player.profileBanner); err != nil {
			log.Fatal(err)
		}

		players = append(players, player)

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// Получаем список ников из таблицы members
	playersFromDB := `SELECT name FROM members;`
	playerRows, err := pool.Query(context.Background(), playersFromDB)
	if err != nil {
		log.Fatalf("Can't get player names: %v\n", err)
	}
	defer playerRows.Close()

	var playerNames []string
	// Помещаем наймена игроков в playerNames
	for playerRows.Next() {
		var name string
		err := playerRows.Scan(&name)
		if err != nil {
			log.Fatalf("Scan error: %v\n", err)
		}
		playerNames = append(playerNames, name)
	}

	totalMembers := gjson.Get(resp, "members.#")
	for i := 0; i < int(totalMembers.Int()); i++ {
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
		// M+ scores
		// mythic_plus_scores_by_season
		// https://raider.io/api/v1/characters/profile?region=eu&realm=howling-fjord&name=%D0%A7%D0%BE%D1%81%D0%BA%D0%B8&fields=mythic_plus_scores_by_season%3Acurrent
		found2 := slices.Contains(playerNames, player.name)

		if found2 {

			// fmt.Println("Found", player.name)
			// Это итеррация по всей полученной таблице members
			for _, p := range players {
				// fmt.Println(p.rank, p.name, p.guild, p.realm, p.race, p.class, p.gender, p.faction, p.achievementPoints, p.profileURL, p.profileBanner)
				if player.name == p.name {
					// time sleep нужыен из за ограничения запросов на стороний API
					time.Sleep(220 * time.Millisecond)
					// mythic plus requests
					// Гет запрос

					// Кодирование имени персонажа в URL-кодированный формат, если не кодировать имя персонажа, то API вернет ошибку, почему не знаю.
					encodedName := url.QueryEscape(player.name)

					// Делаем запрос на API
					url := fmt.Sprintf("https://raider.io/api/v1/characters/profile?region=eu&realm=howling-fjord&name=%s&fields=mythic_plus_scores_by_season:current", encodedName)
					respRio, err := http.Get(url)
					if err != nil {
						log.Fatal(err)
					}
					defer respRio.Body.Close()

					// Читаем данные из запроса
					body, err := io.ReadAll(respRio.Body)
					if err != nil {
						log.Fatal(err)
					}

					// Преобразование в строку
					playerResp := string(body)
					if playerResp == "" {
						log.Fatalf("Failed to fetch player data from API")
					}

					// Достаем текущий рейтинг из gjson.Response
					playerRio := gjson.Get(playerResp, "mythic_plus_scores_by_season.#.scores.all")
					var currRioRating string
					// Конвертируем gjson.Response в string
					for _, s := range playerRio.Array() {
						currRioRating = s.String()
					}

					// fmt.Println("О, привет:" + player.name + " " + p.name)
					if player.rank != p.rank || p.mythic_plus_scores_by_season != currRioRating || player.guild != p.guild || player.realm != p.realm || player.race != p.race || player.gender != p.gender || player.achievementPoints != p.achievementPoints || player.profileURL != p.profileURL || player.profileBanner != p.profileBanner {
						updateQuery := "UPDATE members SET "

						var updates []string

						if player.rank != p.rank {
							updates = append(updates, fmt.Sprintf(`rank = '%s'`, player.rank)) // Использование двойных кавычек для строки и %s для интерполяции
						}
						if player.guild != guild {
							updates = append(updates, fmt.Sprintf("guild = '%s'", player.guild)) // Аналогично
						}
						if player.realm != p.realm {
							updates = append(updates, fmt.Sprintf("realm = '%s'", player.realm)) // Аналогично
						}
						if player.race != p.race && player.race != "Mag'har Orc" {
							raceFix := strings.ReplaceAll(player.race, "'", " ")
							updates = append(updates, fmt.Sprintf("race = '%s'", raceFix)) // Аналогично
						}
						if player.gender != p.gender {
							updates = append(updates, fmt.Sprintf("gender = '%s'", player.gender)) // Аналогично
						}
						if player.achievementPoints != p.achievementPoints {
							updates = append(updates, fmt.Sprintf("achievement_points = '%s'", player.achievementPoints)) // Аналогично
						}
						if player.profileURL != p.profileURL {
							updates = append(updates, fmt.Sprintf("profile_url = '%s'", player.profileURL)) // Аналогично
						}

						if p.mythic_plus_scores_by_season != currRioRating {
							updates = append(updates, fmt.Sprintf("mythic_plus_scores_by_season = '%s'", currRioRating))
						}

						if len(updates) > 0 {
							updateQuery += strings.Join(updates, ", ")
							updateQuery += fmt.Sprintf(` WHERE name = '%s'`, player.name)
							fmt.Println(updateQuery)
							_, err := pool.Exec(ctx, updateQuery)
							if err != nil {
								log.Fatal(err)
							} else {
								logger.Println("Updated: ", player.name, updateQuery)
							}
						}
					}
				}
			}
		} else {
			// fmt.Println(playerJson)
			logger.Println("Player ", name.String(), `not found in players list starting insert`)
			insertObject(player, pool)
		}
	}
	defer fmt.Println("UPDATE PLAYERS DONE")
	defer file.Close()
	defer pool.Close()
}

// Добавляем игрока в базу данных
func insertObject(p Player, pool *pgxpool.Pool) {
	ctx := context.Background()
	// Вставка данных в таблицу members
	_, err := pool.Exec(ctx, `
        INSERT INTO members (rank, name, guild, realm, race, class, gender, faction, achievement_points, profile_url, profile_banner, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP)
    `, p.rank, p.name, p.guild, p.realm, p.race, p.class, p.gender, p.faction, p.achievementPoints, p.profileURL, p.profileBanner)
	if err != nil {
		logger.Println("Can't add player: ", p.name, `to database`, err)
	} else {
		logger.Println("Player ", p.name, `added to database`)
	}
}
