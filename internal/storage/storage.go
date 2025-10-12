package storage

import (
	"book-library/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Обертка над db
type Storage struct {
	db *sqlx.DB
}

// Создание новой
func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// Функция создания таблиц на входе (если их нет)
func (s *Storage) СreateTables() error {
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

func (s *Storage) CreateBook(book *models.Book) error {
	// Валидация обязательных полей
	if err := book.ValidateBook(); err != nil {
		return fmt.Errorf("book validation failed: %w", err)
	}

	// Устанавливаем значения по умолчанию
	book.SetDefaults()

	query := `
		INSERT INTO books (
			title, author, genre, description, status, 
			lent_to, lent_date, room, cabinet, shelf, row, created_at
		) VALUES (
			:title, :author, :genre, :description, :status,
			:lent_to, :lent_date, :room, :cabinet, :shelf, :row, :created_at
		) RETURNING id
	`

	// Выполняем запрос и получаем ID новой книги
	rows, err := s.db.NamedQuery(query, book)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&book.ID)
		if err != nil {
			return fmt.Errorf("failed to get book ID: %w", err)
		}
	}

	return nil
}

func (s *Storage) SeedTestData() error {
	books := []models.Book{
		{
			Title:       "Война и мир",
			Author:      "Лев Толстой", 
			Genre:       "Классика",
			Room:        "Гостиная",
			Cabinet:     1,
			Shelf:       1,
			Row:         1,
		},
		{
			Title:       "Преступление и наказание",
			Author:      "Федор Достоевский",
			Genre:       "Классика", 
			Room:        "Гостиная",
			Cabinet:     1,
			Shelf:       1,
			Row:         2,
		},
	}

	for i := range books {
		err := s.CreateBook(&books[i])
		if err != nil {
			return err
		}
		fmt.Printf("Added book: %s (ID: %d)\n", books[i].Title, books[i].ID)
	}

	return nil
}