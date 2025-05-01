package models

// Song represents basic song information.
type Song struct {
	Filename string `json:"filename"`
	Title    string `json:"title,omitempty"` // Optional title
}
