package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/song"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/google/uuid"
)

// SongServicer defines the interface for song-related business logic.
type SongServicer interface {
	GetLibrary(ctx context.Context) ([]song.SongResponse, error)
	ScanMusicLibrary(ctx context.Context) (*song.MusicScanResult, error)
	GetSongFilePath(songID string) (string, error)
	GetAlbumArtPath(imageFileName string) (string, error)
	// Add other song service methods here (e.g., GetSongByID, UpdateSongMetadata)
}

// songService is the concrete implementation of SongServicer.
type songService struct {
	songDB db.SongDBer
}

// NewSongService creates a new instance of SongService.
func NewSongService(songDB db.SongDBer) SongServicer {
	return &songService{songDB: songDB}
}

// GetLibrary retrieves the song library from the database.
func (s *songService) GetLibrary(ctx context.Context) ([]song.SongResponse, error) {
	modelsSongs, err := s.songDB.GetSongLibrary(ctx)
	if err != nil {
		return nil, err
	}
	var responseSongs []song.SongResponse
	for _, song := range modelsSongs {
		responseSongs = append(responseSongs, s.mapSongToResponse(song))
	}
	return responseSongs, nil
}

func (s *songService) mapSongToResponse(modelSong models.Song) song.SongResponse {
	return song.SongResponse{
		ID:              modelSong.ID,
		Title:           modelSong.Title,
		Artist:          modelSong.Artist,
		Album:           modelSong.Album,
		TrackNumber:     modelSong.TrackNumber,
		Genre:           modelSong.Genre,
		Year:            modelSong.Year,
		DurationSeconds: modelSong.DurationSeconds,
		AudioStreamURL:  fmt.Sprintf("/api/audio/%s", modelSong.ID.String()),
		Filename:        modelSong.Filename,
		CreatedAt:       modelSong.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       modelSong.UpdatedAt.Format(time.RFC3339),
		AlbumArtURL:     fmt.Sprintf("/api/album-art/%s", modelSong.AlbumArtPath),
	}
}

// ScanMusicLibrary triggers the music directory scan and database update.
func (s *songService) ScanMusicLibrary(ctx context.Context) (*song.MusicScanResult, error) {
	log.Println("Starting music library scan...")
	result, err := s.songDB.ScanAndStoreSongs(ctx)
	if err != nil {
		log.Printf("Music library scan failed: %v", err)
		return nil, fmt.Errorf("failed to scan music library: %w", err)
	}
	log.Printf("Music library scan complete. Added: %d, Updated: %d, Removed: %d. Errors: %d",
		result.Added, result.Updated, result.Removed, len(result.Errors))
	if len(result.Errors) > 0 {
		for _, errMsg := range result.Errors {
			log.Printf("Scan Error: %s", errMsg)
		}
	}
	return result, nil
}

// GetSongFilePath retrieves the full file path for a song based on its ID.
func (s *songService) GetSongFilePath(songIDStr string) (string, error) {
	if songIDStr == "" {
		return "", fmt.Errorf("song ID is required")
	}

	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid song ID format: %w", err)
	}

	ctx := context.Background()
	song, err := s.songDB.GetSongByID(ctx, songID)
	if err != nil {
		log.Printf("Service: Song with ID %s not found: %v", songID.String(), err)
		return "", fmt.Errorf("audio file not found")
	}

	musicDir := os.Getenv("MUSIC_DIRECTORY")
	if musicDir == "" {
		return "", fmt.Errorf("MUSIC_DIRECTORY environment variable not set")
	}

	fullPath := filepath.Join(musicDir, song.FilePath)

	// Security check: ensure the resolved path is indeed under MUSIC_DIRECTORY
	absMusicDir, err := filepath.Abs(musicDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute music directory: %w", err)
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute song path: %w", err)
	}

	if !strings.HasPrefix(absFullPath, absMusicDir) {
		return "", fmt.Errorf("attempted path traversal: %s", fullPath)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file not found at path: %s", fullPath)
	}

	return fullPath, nil
}

// GetAlbumArtPath retrieves the full file path for an album art image.
func (s *songService) GetAlbumArtPath(songIDStr string) (string, error) { // Changed parameter name to songIDStr
	if songIDStr == "" {
		return "", fmt.Errorf("song ID is required")
	}

	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid song ID format: %w", err)
	}

	ctx := context.Background() // Or pass context from handler
	song, err := s.songDB.GetSongByID(ctx, songID)
	if err != nil {
		log.Printf("Service: Song with ID %s not found: %v", songID.String(), err)
		return "", fmt.Errorf("album art not found for song ID %s", songIDStr)
	}

	if song.AlbumArtPath == "" {
		return "", fmt.Errorf("no album art path found for song ID %s", songIDStr)
	}

	musicDir := os.Getenv("MUSIC_DIRECTORY")
	if musicDir == "" {
		return "", fmt.Errorf("MUSIC_DIRECTORY environment variable not set")
	}

	// Construct the full path using MUSIC_DIRECTORY and the song's AlbumArtPath
	fullPath := filepath.Join(musicDir, song.AlbumArtPath)

	// Security check: ensure the resolved path is indeed under MUSIC_DIRECTORY
	absMusicDir, err := filepath.Abs(musicDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute music directory: %w", err)
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute album art path: %w", err)
	}

	if !strings.HasPrefix(absFullPath, absMusicDir) {
		return "", fmt.Errorf("attempted path traversal: %s", fullPath)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("album art image not found at path: %s", fullPath)
	}

	return fullPath, nil
}
