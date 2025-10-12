package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage)СreateTables() error {
	fmt.Println("creat tables")

	tables := []string{
		`CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			genre TEXT,
			room TEXT NOT NULL DEFAULT 'Гостиная',
            cabinet INTEGER NOT NULL DEFAULT 1,
			shelf INTEGER NOT NULL DEFAULT 1,
			row INTEGER NOT NULL DEFAULT 1,
			description TEXT,
			status TEXT,
			lent_to TEXT,
			lent_date TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_books_full_location ON books(room, cabinet, shelf, row)`,
		`CREATE INDEX IF NOT EXISTS idx_books_room ON books(room)`,
		`CREATE INDEX IF NOT EXISTS idx_books_author ON books(author)`,
	}

	for _, tableSQL := range tables {
		_, err := s.db.Exec(tableSQL)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	fmt.Println("Tables created successfully!")
	return nil
}
