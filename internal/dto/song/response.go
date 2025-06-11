package song

import "github.com/google/uuid"

// SongResponse defines the response structure for a single song.
type SongResponse struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Artist          string    `json:"artist"`
	Album           string    `json:"album,omitempty"`
	TrackNumber     int       `json:"track_number,omitempty"`
	Genre           string    `json:"genre,omitempty"`
	Year            int       `json:"year,omitempty"`
	DurationSeconds int       `json:"duration_seconds"`
	FilePath        string    `json:"file_path"` // This might be sensitive, consider if it should be exposed
	Filename        string    `json:"filename"`
	CreatedAt       string    `json:"created_at"` // Often sent as string in ISO format
	UpdatedAt       string    `json:"updated_at"`
	AlbumArtPath    string    `json:"album_art_path,omitempty"`
}

// ListSongsResponse defines the response structure for a list of songs.
type ListSongsResponse struct {
	Songs []SongResponse `json:"songs"`
	Total int            `json:"total"` // Useful for pagination
}
