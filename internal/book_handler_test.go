package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type mockService struct {
	getBooksByGenreResponse      Response
	getBooksByGenreError         error
	submitPickUpScheduleResponse PostResponse
	submitPickUpScheduleError    error
}

func (m *mockService) GetBooksByGenreService(ctx context.Context, genre string) (Response, error) {
	return m.getBooksByGenreResponse, m.getBooksByGenreError
}

func (m *mockService) SubmitPickUpScheduleService(schedule PickUpSchedule) (PostResponse, error) {
	return m.submitPickUpScheduleResponse, m.submitPickUpScheduleError
}

func TestBookHandler_GetBooksByGenreHandler(t *testing.T) {
	mockService := &mockService{
		getBooksByGenreResponse: Response{
			Status:    "200 OK",
			IsSuccess: true,
			Message:   "fetch data books successfully!",
			TotalData: 2,
			Data:      []Book{{Title: "Book1"}, {Title: "Book2"}},
		},
		getBooksByGenreError: nil,
	}
	handler := NewHandler(mockService)

	t.Run("PositiveCase", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/books/fiction", nil)
		rec := httptest.NewRecorder()

		router := httprouter.New()
		router.GET("/books/:genre", handler.GetBooksByGenreHandler)

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rec.Code)
		}

		var response Response
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response body: %v", err)
		}

		if response.IsSuccess != true {
			t.Errorf("Expected success, got %v", response.IsSuccess)
		}

		if len(response.Data) != 2 {
			t.Errorf("Expected 2 books, got %d", len(response.Data))
		}
	})

	t.Run("NegativeCase_InternalServerError", func(t *testing.T) {
		mockService.getBooksByGenreError = fmt.Errorf("Internal Server Error")
		req := httptest.NewRequest("GET", "/books/nonexistentgenre", nil)
		rec := httptest.NewRecorder()

		router := httprouter.New()
		router.GET("/books/:genre", handler.GetBooksByGenreHandler)

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %d", rec.Code)
		}
	})
}

func TestBookHandler_SubmitPickUpScheduleHandler(t *testing.T) {
	mockService := &mockService{
		submitPickUpScheduleResponse: PostResponse{
			Status:    "201 CREATED",
			IsSuccess: true,
			Message:   "save new data books successfully!",
			TotalData: 1,
			Data:      PickUpSchedule{},
		},
		submitPickUpScheduleError: nil,
	}
	handler := NewHandler(mockService)

	t.Run("PositiveCase", func(t *testing.T) {
		schedule := PickUpSchedule{
			Genre: "fiction",
			BookInfo: Book{
				Title:         "NewBook",
				Author:        []string{"Author1"},
				EditionNumber: 1,
			},
		}

		reqBody, err := json.Marshal(schedule)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest("POST", "/books/schedule", bytes.NewReader(reqBody))
		rec := httptest.NewRecorder()

		router := httprouter.New()
		router.POST("/books/schedule", handler.SubmitPickUpScheduleHandler)

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code 201, got %d", rec.Code)
		}

		var response PostResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response body: %v", err)
		}

		if response.IsSuccess != true {
			t.Errorf("Expected success, got %v", response.IsSuccess)
		}

		if response.TotalData != 1 {
			t.Errorf("Expected total data 1, got %d", response.TotalData)
		}
	})

	t.Run("NegativeCase_BadRequest", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/books/schedule", bytes.NewReader([]byte("invalid request body")))
		rec := httptest.NewRecorder()

		router := httprouter.New()

		router.POST("/books/schedule", handler.SubmitPickUpScheduleHandler)

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", rec.Code)
		}
	})

}
