package main

import (
	"context"
	"costmart-backend-test/internal"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	// Initialize context
	ctx := context.Background()

	// Initialize http router
	router := httprouter.New()

	// Initialize book module with in-memory storage
	bookRepo := internal.NewInMemoryRepository(ctx)
	bookService := internal.NewService(bookRepo)
	bookHandler := internal.NewHandler(bookService)

	// Define API routes
	router.GET("/books/:genre", bookHandler.GetBooksByGenreHandler)
	router.POST("/books/schedule", bookHandler.SubmitPickUpScheduleHandler)

	// Run the server
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
