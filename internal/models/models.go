package models

import (
	"fmt"
	"time"
)

// Основная структура книги
type Book struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	Genre       string    `json:"genre" db:"genre"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	LentTo      string    `json:"lent_to" db:"lent_to"`
	LentDate    time.Time `json:"lent_date" db:"lent_date"`
	Room        string    `json:"room" db:"room"`
	Cabinet     int       `json:"cabinet" db:"cabinet"`
	Shelf       int       `json:"shelf" db:"shelf"`
	Row         int       `json:"row" db:"row"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// устанавливаем значения по умолчанию
func (b *Book) SetDefaults() {
	if b.Status == "" {
		b.Status = "available"
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
}

//проверяем обязательные поля книги
func (book *Book) ValidateBook() error {
	if book.Title == "" {
		return fmt.Errorf("title is required")
	}
	if book.Author == "" {
		return fmt.Errorf("author is required")
	}
	if book.Room == "" {
		return fmt.Errorf("room is required")
	}
	if book.Cabinet <= 0 {
		return fmt.Errorf("cabinet must be positive")
	}
	if book.Shelf <= 0 {
		return fmt.Errorf("shelf must be positive")
	}
	if book.Row <= 0 {
		return fmt.Errorf("row must be positive")
	}
	return nil
}
