package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

// GetLibrary retrieves the song library from the database and maps them to SongResponse DTOs.
func (s *songService) GetLibrary(ctx context.Context) ([]song.SongResponse, error) {
	songs, err := s.songDB.GetSongLibrary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get song library from DB: %w", err)
	}

	var songResponses []song.SongResponse
	for _, s := range songs {
		songResponses = append(songResponses, mapSongToResponse(&s))
	}
	return songResponses, nil
}

func mapSongToResponse(s *models.Song) song.SongResponse {
	// Example: FilePath "Artist/Album/Song.mp3" -> AlbumArtURL "/api/album-art/Artist/Album/thumb.jpg"
	albumArtURL := ""
	if s.FilePath != "" {
		albumArtURL = "/api/album-art/" + filepath.ToSlash(filepath.Join(filepath.Dir(s.FilePath), "thumb.jpg"))
	}

	return song.SongResponse{
		ID:              s.ID,
		Title:           s.Title,
		Artist:          s.Artist,
		Album:           s.Album,
		TrackNumber:     s.TrackNumber,
		Genre:           s.Genre,
		Year:            s.Year,
		DurationSeconds: s.DurationSeconds,
		AudioStreamURL:  fmt.Sprintf("/api/audio/%s", s.ID.String()),
		Filename:        s.Filename,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		AlbumArtURL:     albumArtURL, // Set the derived URL here
	}
}

func (s *songService) GetAlbumArtPath(songIDStr string) (string, error) {
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
		log.Printf("Service: Song with ID %s not found for album art: %v", songID.String(), err)
		return "", fmt.Errorf("album art not found for song ID %s", songIDStr)
	}

	if song.FilePath == "" { // Check FilePath instead of AlbumArtPath
		return "", fmt.Errorf("song has no file path to derive album art from for ID %s", songIDStr)
	}

	musicDir := os.Getenv("MUSIC_DIRECTORY")
	if musicDir == "" {
		return "", fmt.Errorf("MUSIC_DIRECTORY environment variable not set")
	}
	relativeArtPath := filepath.Join(filepath.Dir(song.FilePath), "thumb.jpg")
	fullPath := filepath.Join(musicDir, relativeArtPath)
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
		// Log this as a warning, not necessarily an error, as not all albums might have art.
		log.Printf("Album art image not found at expected path: %s for song ID %s", fullPath, songIDStr)
		return "", fmt.Errorf("album art image not found for song ID %s", songIDStr) // Still return error to handler
	}

	return fullPath, nil
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
