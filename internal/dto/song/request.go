package song

// CreateSongInput defines the request body for creating a new song.
type CreateSongInput struct {
	Title  string `json:"title" binding:"required"`
	Artist string `json:"artist" binding:"required"`
	Album  string `json:"album"`
	// Add other fields as needed (e.g., Genre, ReleaseYear)
}

// UpdateSongInput defines the request body for updating an existing song.
type UpdateSongInput struct {
	Title  *string `json:"title"` // Use pointers for optional updates
	Artist *string `json:"artist"`
	Album  *string `json:"album"`
}
