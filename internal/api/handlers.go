package api

import (
	"book-library/internal/models"
	"book-library/internal/storage"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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

func GetAllBooks(storage storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		books, err := storage.GetAllBooks(100, 0)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get books",
			})
		}

		return c.JSON(fiber.Map{
			"count": len(books),
			"books": books,
		})
	}
}

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

		// Устанавливаем ID (на всякий случай)
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
