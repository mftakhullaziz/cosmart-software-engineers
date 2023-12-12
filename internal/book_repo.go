package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BookRepository interface {
	GetBooksByGenre(ctx context.Context, genre string) ([]Book, []PickUpSchedule, error)
	SavePickUpSchedule(schedule PickUpSchedule) ([]PickUpSchedule, error)
}

type InMemoryRepository struct {
	ctx                context.Context
	booksWithSchedules map[string]struct {
		Books           []Book
		PickUpSchedules []PickUpSchedule
	}
}

func NewInMemoryRepository(ctx context.Context) *InMemoryRepository {
	return &InMemoryRepository{
		ctx: ctx,
		booksWithSchedules: make(map[string]struct {
			Books           []Book
			PickUpSchedules []PickUpSchedule
		}),
	}
}

func (r *InMemoryRepository) GetBooksByGenre(ctx context.Context, genre string) ([]Book, []PickUpSchedule, error) {
	var newPickUpSchedule []PickUpSchedule

	// Retrieve books and pick-up schedules from the cache
	data, exists := r.booksWithSchedules[genre]
	if exists {
		newPickUpSchedule = append(newPickUpSchedule, data.PickUpSchedules...)
	}

	// Fetch data from API and update cache atomically
	books, err := r.fetchBooksByGenreExternalAPI(ctx, genre)
	if err != nil {
		return nil, nil, err
	}

	return books, newPickUpSchedule, nil
}

func (r *InMemoryRepository) SavePickUpSchedule(schedule PickUpSchedule) ([]PickUpSchedule, error) {
	// Assuming you have the genre information in the schedule
	genre := schedule.Genre

	// Update the cache with the new pick-up schedule
	if data, exists := r.booksWithSchedules[genre]; exists {
		data.PickUpSchedules = append(data.PickUpSchedules, schedule)
		r.booksWithSchedules[genre] = data
	} else {
		r.booksWithSchedules[genre] = struct {
			Books           []Book
			PickUpSchedules []PickUpSchedule
		}{
			PickUpSchedules: []PickUpSchedule{schedule},
		}
	}

	return r.booksWithSchedules[genre].PickUpSchedules, nil
}

func (r *InMemoryRepository) fetchBooksByGenreExternalAPI(ctx context.Context, genre string) ([]Book, error) {
	// Build the URL with the specified genre
	url := fmt.Sprintf("https://openlibrary.org/subjects/%s.json", genre)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err.Error())
		}
	}(response.Body)

	// Print response
	fmt.Println(response.Body)

	// Check if the response status code is not 200 OK
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", response.StatusCode)
	}

	// Read and parse the response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	// Extract relevant book information from the response
	books, err := constructOfResponse(data)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func constructOfResponse(data map[string]interface{}) ([]Book, error) {
	var books []Book
	var authorsList []string

	if subjects, ok := data["works"].([]interface{}); ok {
		for _, subject := range subjects {
			if work, ok := subject.(map[string]interface{}); ok {
				// Extracting the title
				title := ""
				if titleValue, exists := work["title"].(string); exists {
					title = titleValue
				}

				// Extracting the authors
				var authors []string
				if authorsArray, exists := work["authors"].([]interface{}); exists {
					for _, author := range authorsArray {
						if authorMap, isMap := author.(map[string]interface{}); isMap {
							if authorName, hasName := authorMap["name"].(string); hasName {
								authors = append(authors, authorName)
								authorsList = append(authorsList, authorName)
							}
						}
					}
				}

				// Create a Book instance and append it to the books slice
				book := Book{
					Title:         title,
					Author:        authors,
					EditionNumber: int(work["edition_count"].(float64)),
				}
				books = append(books, book)
			}
		}
	}

	return books, nil
}
