package structs

import "time"

type Entry struct {
	Title     string    `json:"title" xml:"title"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at" xml:"published"`
}
