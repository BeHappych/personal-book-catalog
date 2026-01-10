package api

import (
	"book-library/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type (
	Api struct {
		fiber *fiber.App
	}
)

func NewApi(storage storage.Storage) *Api {
	api := Api{fiber: fiber.New(fiber.Config{AppName: "Book Library"})}

	api.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	api.fiber.Static("/", "./static")
	api.fiber.Get("/", func(c *fiber.Ctx) error {
		books, _ := storage.GetAllBooks(10, 0)
		return c.Render("index", fiber.Map{
			"Books": books,
		})
	})

	api.fiber.Post("/api/books", AddBook(storage))
	api.fiber.Get("/api/books/:id", GetBookByID(storage))
	api.fiber.Get("/api/books", GetAllBooks(storage))
	api.fiber.Delete("/api/books/:id", DeleteBookByID(storage))
	api.fiber.Put("/api/books/:id", UpdateBookByID(storage))
	api.fiber.Post("/api/books/:id/lend", LendBook(storage))
    api.fiber.Post("/api/books/:id/return", ReturnBook(storage))

	return &api
}

func (api *Api) Start(addr string) {

	go func() {
		if err := api.fiber.Listen(addr); err != nil {
			panic(err)
		}
	}()
}
