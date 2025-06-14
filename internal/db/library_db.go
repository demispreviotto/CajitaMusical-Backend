package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
)

func GetSongLibrary(ctx context.Context) ([]*models.Song, error) {
	musicDir := os.Getenv("MUSIC_DIRECTORY")

	if musicDir == "" {
		return nil, fmt.Errorf("MUSIC_DIRECTORY environment variable not set")
	}

	var songs []*models.Song
	err := filepath.Walk(musicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".mp3" { // For now, only consider MP3 files
			filename := filepath.Base(path)
			songs = append(songs, &models.Song{Filename: filename})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read song library: %w", err)
	}

	return songs, nil
}

// sanitizeFilename es una función auxiliar para crear nombres de archivo seguros
func sanitizeFilename(name string) string {
	// Reemplaza caracteres no permitidos con un guion bajo o similar. Quita espacios al inicio/final
	sanitized := strings.TrimSpace(name)
	sanitized = strings.ReplaceAll(sanitized, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, "\\", "-")
	sanitized = strings.ReplaceAll(sanitized, ":", "-")
	sanitized = strings.ReplaceAll(sanitized, "*", "-")
	sanitized = strings.ReplaceAll(sanitized, "?", "")
	sanitized = strings.ReplaceAll(sanitized, "\"", "'")
	sanitized = strings.ReplaceAll(sanitized, "<", "")
	sanitized = strings.ReplaceAll(sanitized, ">", "")
	sanitized = strings.ReplaceAll(sanitized, "|", "")
	return sanitized
}
