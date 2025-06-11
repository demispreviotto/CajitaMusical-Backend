package db

import (
	"context"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
)

// SongDBer defines the interface for song database operations.
type SongDBer interface {
	GetSongLibrary(ctx context.Context) ([]models.Song, error)
	// Add other song-related DB methods here (e.g., CreateSong, GetSongByID, UpdateSong, DeleteSong)
}

// songDB is the concrete implementation of SongDBer.
type songDB struct{}

// NewSongDB creates a new instance of SongDB.
func NewSongDB() SongDBer {
	return &songDB{}
}

// GetSongLibrary retrieves all songs from the database.
func (sdb *songDB) GetSongLibrary(ctx context.Context) ([]models.Song, error) {
	var songs []models.Song
	err := DB.WithContext(ctx).Find(&songs).Error
	return songs, err
}

// You'd add other song DB functions here as you implement them, e.g.:
/*
func (sdb *songDB) CreateSong(ctx context.Context, song *models.Song) error {
	return DB.WithContext(ctx).Create(song).Error
}

func (sdb *songDB) GetSongByID(ctx context.Context, songID uuid.UUID) (*models.Song, error) {
	var song models.Song
	err := DB.WithContext(ctx).Where("id = ?", songID).First(&song).Error
	return &song, err
}
*/
