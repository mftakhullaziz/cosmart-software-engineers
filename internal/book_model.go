package internal

type Book struct {
	Title         string   `json:"title"`
	Author        []string `json:"author"`
	EditionNumber int      `json:"edition_number"`
}

type PickUpSchedule struct {
	BookInfo   Book   `json:"book_info"`
	PickUpDate string `json:"pick_up_date"`
	Genre      string `json:"genre"`
}
