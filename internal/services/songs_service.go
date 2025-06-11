package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/song"
)

// SongServicer define la interfaz para las operaciones del servicio de canciones.
//
//go:generate mockgen -source=song_service.go -destination=mocks/mock_song_service.go
type SongServicer interface {
	GetLibrary(ctx context.Context) ([]song.SongResponse, error)
	GetSongFilePath(filename string) (string, error)
	// Add other song-related service methods
}

// songService es la implementación concreta de SongServicer.
type songService struct {
	songDB   db.SongDBer
	musicDir string
}

// NewSongService crea una nueva instancia de SongService.
func NewSongService(songDB db.SongDBer, musicDir string) SongServicer {
	return &songService{songDB: songDB, musicDir: musicDir}
}

// GetLibrary retrieves the song library and maps it to DTOs.
func (s *songService) GetLibrary(ctx context.Context) ([]song.SongResponse, error) {
	modelsSongs, err := s.songDB.GetSongLibrary(ctx)
	if err != nil {
		log.Printf("Service: Failed to fetch song library from DB: %v", err)
		return nil, errors.New("failed to retrieve song library")
	}

	var responseSongs []song.SongResponse
	for _, modelSong := range modelsSongs {
		responseSongs = append(responseSongs, song.SongResponse{
			ID:              modelSong.ID,
			Title:           modelSong.Title,
			Artist:          modelSong.Artist,
			Album:           modelSong.Album,
			TrackNumber:     modelSong.TrackNumber,
			Genre:           modelSong.Genre,
			Year:            modelSong.Year,
			DurationSeconds: modelSong.DurationSeconds,
			FilePath:        modelSong.FilePath,
			Filename:        modelSong.Filename,
			CreatedAt:       modelSong.CreatedAt.Format(time.RFC3339), // Standardized format
			UpdatedAt:       modelSong.UpdatedAt.Format(time.RFC3339), // Standardized format
			AlbumArtPath:    modelSong.AlbumArtPath,
		})
	}

	return responseSongs, nil
}

// GetSongFilePath constructs and validates the full path to a song file.
func (s *songService) GetSongFilePath(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("filename is required")
	}

	filePath := filepath.Join(s.musicDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("audio file not found")
	}

	return filePath, nil
}
