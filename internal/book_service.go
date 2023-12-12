package internal

import (
	"context"
	"fmt"
)

type BookService interface {
	GetBooksByGenreService(ctx context.Context, genre string) (Response, error)
	SubmitPickUpScheduleService(schedule PickUpSchedule) (PostResponse, error)
}

type Response struct {
	Status    string `json:"status"`
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	TotalData int    `json:"total_data"`
	Data      []Book `json:"data"`
}

type PostResponse struct {
	Status    string         `json:"status"`
	IsSuccess bool           `json:"is_success"`
	Message   string         `json:"message"`
	TotalData int            `json:"total_data"`
	Data      PickUpSchedule `json:"data"`
}

type bookService struct {
	repository BookRepository
}

func NewService(repository BookRepository) BookService {
	return &bookService{
		repository: repository,
	}
}

func (s *bookService) GetBooksByGenreService(ctx context.Context, genre string) (Response, error) {
	books, pickUpSchedules, err := s.repository.GetBooksByGenre(ctx, genre)
	if err != nil {
		return Response{
			Status:    "500 Internal Server Error",
			IsSuccess: false,
			Message:   fmt.Sprintf("failed to fetch data books: %v", err),
			Data:      []Book{},
			TotalData: 0,
		}, err
	}

	totalData := len(books) + len(pickUpSchedules)

	// Prepare the response structure
	response := Response{
		Status:    "200 OK",
		IsSuccess: true,
		Message:   "fetch data books successfully!",
		TotalData: totalData,
		Data:      nil, // Initialize with nil slice to avoid null in JSON response
	}

	// Append books to the response
	response.Data = append(response.Data, books...)

	// Append pick-up schedules to the response
	for _, schedule := range pickUpSchedules {
		response.Data = append(response.Data, Book{
			Title:         schedule.BookInfo.Title,
			Author:        schedule.BookInfo.Author,
			EditionNumber: schedule.BookInfo.EditionNumber,
		})
	}

	return response, nil
}

func (s *bookService) SubmitPickUpScheduleService(schedule PickUpSchedule) (PostResponse, error) {
	pickUpSchedule, err := s.repository.SavePickUpSchedule(schedule)
	if err != nil {
		return PostResponse{
			Status:    "500 Internal Server Error",
			IsSuccess: false,
			Message:   fmt.Sprintf("failed to save new data books: %v", err),
			Data:      PickUpSchedule{},
			TotalData: 0,
		}, err
	}

	response := PostResponse{
		Status:    "201 CREATED",
		IsSuccess: true,
		Message:   "save new data books successfully!",
		TotalData: len(pickUpSchedule),
		Data:      pickUpSchedule[0],
	}

	return response, nil
}
