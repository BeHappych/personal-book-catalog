package main

import (
	"book-library/internal/storage"
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	initConfigs()

	db := initDB()
	defer db.Close()

	storage := storage.New(db)

	storage.SeedTestData()

	err := storage.СreateTables()
	if err != nil {
		panic(fmt.Errorf("failed to create tables: %w", err))
	}
}

func initConfigs() {
	fmt.Println("init configs")

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}

		//добавляем функционал чтения конфигурации из переменных окружения
		//будет использоваться если удалить/переименовать config.yaml
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	}
}

func initDB() *sqlx.DB {
	url := getDatabaseURL()

	fmt.Println("init database")

	db, err := sqlx.Connect("postgres", url.String())
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	fmt.Println("Database connected successfully!")
	return db
}

func getDatabaseURL() *url.URL {
	return &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			viper.GetString("db.user"),
			viper.GetString("db.password"),
		),
		Host: fmt.Sprintf("%s:%s",
			viper.GetString("db.address"),
			viper.GetString("db.port"),
		),
		Path:     viper.GetString("db.database"),
		RawQuery: "sslmode=disable",
	}
}
