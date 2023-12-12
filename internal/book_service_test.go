package internal

import (
	"context"
	"fmt"
	"testing"
)

type mockRepository struct {
	getBooksByGenreResponse    []Book
	getPickUpSchedulesResponse []PickUpSchedule
	getBooksByGenreError       error
	savePickUpScheduleResponse []PickUpSchedule
	savePickUpScheduleError    error
}

func (m *mockRepository) GetBooksByGenre(ctx context.Context, genre string) ([]Book, []PickUpSchedule, error) {
	return m.getBooksByGenreResponse, m.getPickUpSchedulesResponse, m.getBooksByGenreError
}

func (m *mockRepository) SavePickUpSchedule(schedule PickUpSchedule) ([]PickUpSchedule, error) {
	return m.savePickUpScheduleResponse, m.savePickUpScheduleError
}

func TestBookService_GetBooksByGenreService(t *testing.T) {
	// Positive case: Books exist in the cache
	mockRepo := &mockRepository{
		getBooksByGenreResponse: []Book{{Title: "MockBook1"}, {Title: "MockBook2"}},
		getPickUpSchedulesResponse: []PickUpSchedule{
			{
				Genre: "fiction",
				BookInfo: Book{
					Title:         "PickUpBook1",
					Author:        []string{"Author1"},
					EditionNumber: 1,
				},
			},
			{
				Genre: "fiction",
				BookInfo: Book{
					Title:         "PickUpBook2",
					Author:        []string{"Author2"},
					EditionNumber: 2,
				},
			},
		},
	}

	service := NewService(mockRepo)

	t.Run("PositiveCase_CacheHit_WithPickUpSchedules", func(t *testing.T) {
		// Perform the test
		response, err := service.GetBooksByGenreService(context.Background(), "fiction")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the expected outcome
		if response.IsSuccess != true {
			t.Errorf("Expected success, got %v", response.IsSuccess)
		}

		if len(response.Data) != 4 {
			t.Errorf("Expected 4 books (2 from cache and 2 from pick-up schedules), got %d", len(response.Data))
		}
	})

	t.Run("PositiveCase_CacheHit", func(t *testing.T) {
		// Perform the test
		response, err := service.GetBooksByGenreService(context.Background(), "fiction")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the expected outcome
		if response.IsSuccess != true {
			t.Errorf("Expected success, got %v", response.IsSuccess)
		}

		if len(response.Data) != 4 {
			t.Errorf("Expected 4 books, got %d", len(response.Data))
		}
	})

	// Negative case: Books do not exist in the cache, API request fails
	mockRepo.getBooksByGenreError = fmt.Errorf("API request failed")
	t.Run("NegativeCase_CacheMiss_APIFailure", func(t *testing.T) {
		// Perform the test
		response, err := service.GetBooksByGenreService(context.Background(), "nonexistentgenre")
		if err == nil {
			t.Error("Expected error, but got nil")
		}

		// Add your assertions here based on the expected error
		if response.IsSuccess != false {
			t.Errorf("Expected failure, got %v", response.IsSuccess)
		}
	})
}

func TestBookService_SubmitPickUpScheduleService(t *testing.T) {
	mockRepo := &mockRepository{
		savePickUpScheduleResponse: []PickUpSchedule{{BookInfo: Book{Title: "MockBook"}}},
	}

	service := NewService(mockRepo)

	t.Run("PositiveCase", func(t *testing.T) {
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
		response, err := service.SubmitPickUpScheduleService(schedule)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Add your assertions here based on the expected outcome
		if response.IsSuccess != true {
			t.Errorf("Expected success, got %v", response.IsSuccess)
		}

		if len(response.Data.BookInfo.Title) == 0 {
			t.Errorf("Expected non-empty title, got empty")
		}
	})

	t.Run("NegativeCase_SaveFailure", func(t *testing.T) {
		// Set an error for saving pick-up schedule
		mockRepo.savePickUpScheduleError = fmt.Errorf("Failed to save schedule")

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
		response, err := service.SubmitPickUpScheduleService(schedule)
		if err == nil {
			t.Error("Expected error, but got nil")
		}

		// Add your assertions here based on the expected error
		if response.IsSuccess != false {
			t.Errorf("Expected failure, got %v", response.IsSuccess)
		}
	})
}
