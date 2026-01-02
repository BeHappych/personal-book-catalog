package api

import (
	"book-library/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type (
	Api struct {
		fiber *fiber.App
	}
)

func NewApi(storage storage.Storage) *Api {
	api := Api{fiber: fiber.New(fiber.Config{})}

	api.fiber.Get("/api/book/:id", GetBookByID(storage))
	api.fiber.Get("/api/books", GetAllBooks(storage))
	api.fiber.Delete("/api/books/:id", DeleteBookByID(storage))
	api.fiber.Put("/api/books/:id", UpdateBookByID(storage))

	return &api
}

func (api *Api) Start(addr string) {

	go func() {
		if err := api.fiber.Listen(addr); err != nil {
			panic(err)
		}
	}()
}
