package models

import "time"

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
