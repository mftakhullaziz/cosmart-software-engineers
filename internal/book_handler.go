package internal

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type BookHandler interface {
	GetBooksByGenreHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	SubmitPickUpScheduleHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type bookHandler struct {
	service BookService
}

func NewHandler(service BookService) BookHandler {
	return &bookHandler{
		service: service,
	}
}

func (h *bookHandler) GetBooksByGenreHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	genre := params.ByName("genre")
	books, err := h.service.GetBooksByGenreService(r.Context(), genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func (h *bookHandler) SubmitPickUpScheduleHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var schedule PickUpSchedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err.Error())
		}
	}(r.Body)

	pickUpSchedule, err := h.service.SubmitPickUpScheduleService(schedule)
	response, err := json.Marshal(pickUpSchedule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		return
	}
}
