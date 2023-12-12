package internal

import (
	"encoding/json"
	"reflect"
	"testing"
)

// equalSlices returns true if two slices are equal, false otherwise.
func equalSlices(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}

func TestBookSerialization(t *testing.T) {
	// Create a sample Book instance
	book := Book{
		Title:         "Sample Book",
		Author:        []string{"Author1", "Author2"},
		EditionNumber: 1,
	}

	// Serialize the Book instance to JSON
	jsonData, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal Book to JSON: %v", err)
	}

	// Deserialize the JSON back to a Book instance
	var deserializedBook Book
	err = json.Unmarshal(jsonData, &deserializedBook)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to Book: %v", err)
	}

	// Check if the original and deserialized Book instances are equal
	if deserializedBook.Title != book.Title ||
		deserializedBook.EditionNumber != book.EditionNumber ||
		!equalSlices(deserializedBook.Author, book.Author) {
		t.Error("Original and deserialized Book instances are not equal")
	}
}

func TestPickUpScheduleSerialization(t *testing.T) {
	// Create a sample PickUpSchedule instance
	pickUpSchedule := PickUpSchedule{
		BookInfo: Book{
			Title:         "Sample Book",
			Author:        []string{"Author1", "Author2"},
			EditionNumber: 1,
		},
		PickUpDate: "2023-12-31",
		Genre:      "Fiction",
	}

	// Serialize the PickUpSchedule instance to JSON
	jsonData, err := json.Marshal(pickUpSchedule)
	if err != nil {
		t.Fatalf("Failed to marshal PickUpSchedule to JSON: %v", err)
	}

	// Deserialize the JSON back to a PickUpSchedule instance
	var deserializedPickUpSchedule PickUpSchedule
	err = json.Unmarshal(jsonData, &deserializedPickUpSchedule)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to PickUpSchedule: %v", err)
	}

	// Check if the original and deserialized PickUpSchedule instances are equal
	if deserializedPickUpSchedule.PickUpDate != pickUpSchedule.PickUpDate ||
		deserializedPickUpSchedule.Genre != pickUpSchedule.Genre ||
		!reflect.DeepEqual(deserializedPickUpSchedule.BookInfo, pickUpSchedule.BookInfo) {
		t.Error("Original and deserialized PickUpSchedule instances are not equal")
	}
}
