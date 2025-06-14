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
	AudioStreamURL  string    `json:"audio_stream_url"`
	Filename        string    `json:"filename"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	AlbumArtURL     string    `json:"album_art_url,omitempty"`
}

// ListSongsResponse defines the response structure for a list of songs.
type ListSongsResponse struct {
	Songs []SongResponse `json:"songs"`
	Total int            `json:"total"` // Useful for pagination
}

// MusicScanResult defines the result of a music scan operation.
type MusicScanResult struct {
	Added   int
	Updated int
	Removed int
	Errors  []string
}
