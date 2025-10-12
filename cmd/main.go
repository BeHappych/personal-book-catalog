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

// func createTables(db *sqlx.DB) {
// 	fmt.Println("creat tables")

// 	tables := []string{
// 		`CREATE TABLE IF NOT EXISTS books (
// 			id SERIAL PRIMARY KEY,
// 			title TEXT NOT NULL,
// 			author TEXT NOT NULL,
// 			genre TEXT,
// 			room TEXT NOT NULL DEFAULT 'Гостиная',
//             cabinet INTEGER NOT NULL DEFAULT 1,
// 			shelf INTEGER NOT NULL DEFAULT 1,
// 			row INTEGER NOT NULL DEFAULT 1,
// 			description TEXT,
// 			status TEXT,
// 			lent_to TEXT,
// 			lent_date TIMESTAMP,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 		)`,
//         `CREATE INDEX IF NOT EXISTS idx_books_full_location ON books(room, cabinet, shelf, row)`,
//         `CREATE INDEX IF NOT EXISTS idx_books_room ON books(room)`,
//         `CREATE INDEX IF NOT EXISTS idx_books_author ON books(author)`,
// 	}

// 	for _, tableSQL := range tables {
// 		_, err := db.Exec(tableSQL)
// 		if err != nil {
// 			panic(fmt.Errorf("failed to create table: %w", err))
// 		}
// 	}

// 	fmt.Println("Tables created successfully!")
// }

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
