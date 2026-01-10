package api

import (
	"book-library/internal/models"
	"book-library/internal/storage"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Функция добавления книги
func AddBook(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var book models.Book

		// Парсим JSON из тела запроса
		if err := c.BodyParser(&book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Создаем книгу в базе
		err := storage.CreateBook(&book)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Failed to create book",
				"details": err.Error(),
			})
		}

		// Возвращаем созданную книгу с ID
		return c.Status(fiber.StatusCreated).JSON(book)
	}
}

// Функция получения книги по id
func GetBookByID(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).Send(nil)
		}

		book, err := storage.GetBookByID(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).Send(nil)
		}

		return c.JSON(book)
	}
}

// Функция получения всех книг
func GetAllBooks(storage storage.Storage) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Получаем параметры фильтрации из query string
        title := c.Query("title")
        author := c.Query("author")
		genre := c.Query("genre")
        status := c.Query("status")
        
        // Получаем книги с фильтрами
        books, err := storage.GetBooksWithFilters(models.BookFilters{
            Title:  title,
            Author: author,
			Genre:  genre,
            Status: status,
            Limit:  100, 
            Offset: 0,
        })
        
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to get books",
                "details": err.Error(),
            })
        }
        
        // Если books == nil, возвращаем пустой массив
        if books == nil {
            books = []models.Book{}
        }
        
        // Возвращаем массив книг
        return c.JSON(books)
    }
}

// Функция удаления книги по id
func DeleteBookByID(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid book ID",
			})
		}

		if id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Book ID must be positive",
			})
		}

		// Выполняем удаление
		err = storage.DeleteBook(id)
		if err != nil {
			// Проверяем тип ошибки
			if strings.Contains(err.Error(), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Book not found",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to delete book",
				"details": err.Error(),
			})
		}

		// Возвращаем успешный ответ
		return c.JSON(fiber.Map{
			"message": "Book deleted successfully",
			"id":      id,
		})
	}
}

// Функция обновления книги по id
func UpdateBookByID(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid book ID",
			})
		}

		var updateData map[string]interface{}
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Получаем текущую книгу из БД
		book, err := storage.GetBookByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}

		// Обновляем только переданные поля (частичное обновление)
		updateBookFields(book, updateData)

		// Устанавливаем ID
		book.ID = id

		// Выполняем обновление
		err = storage.UpdateBook(book)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update book",
				"details": err.Error(),
			})
		}

		// Возвращаем обновленную книгу
		return c.JSON(fiber.Map{
			"message": "Book updated successfully",
			"book":    book,
		})
	}
}

// Вспомогательная функция для частичного обновления полей
func updateBookFields(book *models.Book, data map[string]interface{}) {
	if title, ok := data["title"].(string); ok && title != "" {
		book.Title = title
	}
	if author, ok := data["author"].(string); ok && author != "" {
		book.Author = author
	}
	if genre, ok := data["genre"].(string); ok {
		book.Genre = genre
	}
	if description, ok := data["description"].(string); ok {
		book.Description = description
	}
	if status, ok := data["status"].(string); ok && status != "" {
		book.Status = status
	}
	if lentTo, ok := data["lent_to"].(string); ok {
		book.LentTo = lentTo
	}
	if room, ok := data["room"].(string); ok && room != "" {
		book.Room = room
	}
	if cabinet, ok := data["cabinet"].(float64); ok && cabinet > 0 {
		book.Cabinet = int(cabinet)
	}
	if shelf, ok := data["shelf"].(float64); ok && shelf > 0 {
		book.Shelf = int(shelf)
	}
	if row, ok := data["row"].(float64); ok && row > 0 {
		book.Row = int(row)
	}
}

// Функция выдачи книги
func LendBook(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем ID книги
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid book ID",
			})
		}

		// Получаем данные из запроса
		var request struct {
			LentTo string `json:"lent_to"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Проверяем обязательное поле
		if request.LentTo == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Field 'lent_to' is required",
			})
		}

		// Получаем книгу из базы
		book, err := storage.GetBookByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}

		// Проверяем, что книга доступна для выдачи
		if book.Status != "available" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Book is already %s", book.Status),
			})
		}

		// Обновляем данные книги
		book.Status = "lent"
		book.LentTo = request.LentTo
		book.LentDate = time.Now()

		// Сохраняем изменения
		err = storage.UpdateBook(book)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update book",
				"details": err.Error(),
			})
		}

		// Возвращаем обновленную книгу
		return c.JSON(book)
	}
}

// Функция возвращения книги
func ReturnBook(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем ID книги
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid book ID",
			})
		}

		// Получаем книгу из базы
		book, err := storage.GetBookByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}

		// Проверяем, что книга выдана
		if book.Status != "lent" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Book is not lent (current status: %s)", book.Status),
			})
		}

		// Возвращаем книгу
		book.Status = "available"
		book.LentTo = ""
		book.LentDate = time.Time{} 

		// Сохраняем изменения
		err = storage.UpdateBook(book)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update book",
				"details": err.Error(),
			})
		}

		return c.JSON(book)
	}
}
