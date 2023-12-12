package internal

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestInMemoryRepository_GetBooksByGenre(t *testing.T) {
	// Create a new httpmock instance
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register a mock response for the API request
	genre := "fiction"
	mockURL := fmt.Sprintf("https://openlibrary.org/subjects/%s.json", genre)
	mockResponseBody := `{"works": [{"title": "MockBook", "authors": [{"key": "authors/001AAS", "name": "authors"}], "edition_count": 1}]}`
	httpmock.RegisterResponder("GET", mockURL, httpmock.NewStringResponder(200, mockResponseBody))

	// Initialize the repository
	ctx := context.Background()
	repo := NewInMemoryRepository(ctx)

	t.Run("PositiveCase", func(t *testing.T) {
		// Perform the test
		books, pickUpSchedules, err := repo.GetBooksByGenre(ctx, genre)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the mocked response
		if len(books) != 1 {
			t.Errorf("Expected 1 book, got %d", len(books))
		}

		if len(pickUpSchedules) != 0 {
			t.Errorf("Expected 0 pick-up schedules, got %d", len(pickUpSchedules))
		}
	})

	t.Run("NegativeCase", func(t *testing.T) {
		// Register a mock response for the non-existent genre
		nonExistentGenre := "non existent genre"
		mockNonExistentURL := fmt.Sprintf("https://openlibrary.org/subjects/%s.json", nonExistentGenre)
		httpmock.RegisterResponder("GET", mockNonExistentURL, httpmock.NewStringResponder(404, "Not Found"))

		// Perform the test
		_, _, err := repo.GetBooksByGenre(ctx, nonExistentGenre)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
	})

	t.Run("PositiveCase_CacheExist", func(t *testing.T) {
		// Add data to the cache
		repo.booksWithSchedules[genre] = struct {
			Books           []Book
			PickUpSchedules []PickUpSchedule
		}{
			Books: []Book{},
			PickUpSchedules: []PickUpSchedule{{BookInfo: Book{
				"Book Cache",
				[]string{"Author1", "Author2"},
				1,
			}, Genre: genre}},
		}

		// Perform the test
		books, pickUpSchedules, err := repo.GetBooksByGenre(ctx, genre)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the cached response
		if len(books) != 1 {
			t.Errorf("Expected 1 book, got %d", len(books))
		}

		if len(pickUpSchedules) != 1 {
			t.Errorf("Expected 1 pick-up schedule, got %d", len(pickUpSchedules))
		}
	})

	t.Run("PositiveCase_CacheNotExist", func(t *testing.T) {
		// Clear the cache to simulate cache not containing any data
		repo.booksWithSchedules = make(map[string]struct {
			Books           []Book
			PickUpSchedules []PickUpSchedule
		})

		// Perform the test
		books, pickUpSchedules, err := repo.GetBooksByGenre(ctx, genre)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the response from the API
		if len(books) != 1 {
			t.Errorf("Expected 1 book, got %d", len(books))
		}

		if len(pickUpSchedules) != 0 {
			t.Errorf("Expected 0 pick-up schedules, got %d", len(pickUpSchedules))
		}
	})
}

func TestInMemoryRepository_SavePickUpSchedule(t *testing.T) {
	// Initialize the repository
	ctx := context.Background()
	repo := NewInMemoryRepository(ctx)

	t.Run("SavePickUpScheduleTest", func(t *testing.T) {
		// Create a pick-up schedule
		schedule := PickUpSchedule{
			Genre: "fiction",
			BookInfo: Book{
				Title:         "TestBook",
				Author:        []string{"TestAuthor"},
				EditionNumber: 1,
			},
		}

		// Perform the test
		pickUpSchedules, err := repo.SavePickUpSchedule(schedule)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the expected outcome
		if len(pickUpSchedules) != 1 {
			t.Errorf("Expected 1 pick-up schedule, got %d", len(pickUpSchedules))
		}
	})

}
