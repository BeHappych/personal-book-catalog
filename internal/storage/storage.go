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

// Получение книги по ID
func (s *Storage) GetBookByID(id int) (*models.Book, error) {
    var book models.Book
    query := `SELECT * FROM books WHERE id = $1`
    err := s.db.Get(&book, query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }
    return &book, nil
}

// Получение всех книг
func (s *Storage) GetAllBooks(limit, offset int) ([]models.Book, error) {
    var books []models.Book
    query := `SELECT * FROM books ORDER BY id LIMIT $1 OFFSET $2`
    err := s.db.Select(&books, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to get books: %w", err)
    }
    return books, nil
}

// Обновление книги
func (s *Storage) UpdateBook(book *models.Book) error {
    if err := book.ValidateBook(); err != nil {
        return fmt.Errorf("book validation failed: %w", err)
    }
    
    query := `
        UPDATE books SET
            title = :title,
            author = :author,
            genre = :genre,
            description = :description,
            status = :status,
            lent_to = :lent_to,
            lent_date = :lent_date,
            room = :room,
            cabinet = :cabinet,
            shelf = :shelf,
            row = :row
        WHERE id = :id
    `
    
    result, err := s.db.NamedExec(query, book)
    if err != nil {
        return fmt.Errorf("failed to update book: %w", err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("book with id %d not found", book.ID)
    }
    
    return nil
}

// Удаление книги
func (s *Storage) DeleteBook(id int) error {
    query := `DELETE FROM books WHERE id = $1`
    result, err := s.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete book: %w", err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("book with id %d not found", id)
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