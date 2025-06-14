package db

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/song"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/dhowden/tag"
	"github.com/google/uuid"
)

// SongDBer defines the interface for song database operations.
type SongDBer interface {
	GetSongLibrary(ctx context.Context) ([]models.Song, error)
	CreateSong(ctx context.Context, song *models.Song) error
	UpdateSong(ctx context.Context, song *models.Song) error
	GetSongByFilePath(ctx context.Context, filePath string) (*models.Song, error)
	DeleteSong(ctx context.Context, songID uuid.UUID) error
	ScanAndStoreSongs(ctx context.Context) (*song.MusicScanResult, error)
	GetSongByID(ctx context.Context, songID uuid.UUID) (*models.Song, error)
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

// CreateSong adds a new song to the database.
func (sdb *songDB) CreateSong(ctx context.Context, song *models.Song) error {
	return DB.WithContext(ctx).Create(song).Error
}

// UpdateSong updates an existing song in the database.
func (sdb *songDB) UpdateSong(ctx context.Context, song *models.Song) error {
	return DB.WithContext(ctx).Save(song).Error
}

// GetSongByFilePath retrieves a song by its file path from the database.
func (sdb *songDB) GetSongByFilePath(ctx context.Context, filePath string) (*models.Song, error) {
	var song models.Song
	err := DB.WithContext(ctx).Where("file_path = ?", filePath).First(&song).Error
	return &song, err
}

// DeleteSong removes a song from the database by its ID.
func (sdb *songDB) DeleteSong(ctx context.Context, songID uuid.UUID) error {
	return DB.WithContext(ctx).Delete(&models.Song{}, songID).Error
}

// GetSongByID retrieves a song by its ID from the database.
func (sdb *songDB) GetSongByID(ctx context.Context, songID uuid.UUID) (*models.Song, error) {
	var song models.Song
	if err := DB.WithContext(ctx).First(&song, songID).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

// ScanAndStoreSongs escanea el directorio de música y actualiza la base de datos.
func (sdb *songDB) ScanAndStoreSongs(ctx context.Context) (*song.MusicScanResult, error) {
	musicDir := os.Getenv("MUSIC_DIRECTORY")
	result := &song.MusicScanResult{}

	if musicDir == "" {
		return nil, fmt.Errorf("MUSIC_DIRECTORY environment variable not set")
	}

	var existingSongs []models.Song
	if err := DB.WithContext(ctx).Find(&existingSongs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch existing songs from DB: %w", err)
	}
	existingSongMap := make(map[string]models.Song)
	for _, s := range existingSongs {
		existingSongMap[s.FilePath] = s
	}

	newlyScannedPaths := make(map[string]bool)

	err := filepath.Walk(musicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error accessing path %s: %v", path, err))
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".mp3" && ext != ".flac" && ext != ".m4a" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file %s: %v", path, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error opening file %s: %v", path, err))
			return nil
		}
		defer f.Close()

		t, err := tag.ReadFrom(f)
		if err != nil {
			log.Printf("Error reading tags from %s: %v", path, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error reading tags from %s: %v", path, err))
			return nil
		}

		normalizedMusicDir := musicDir
		if !strings.HasSuffix(normalizedMusicDir, string(os.PathSeparator)) {
			normalizedMusicDir += string(os.PathSeparator)
		}
		relativeFilePath := strings.TrimPrefix(path, normalizedMusicDir)
		relativeFilePath = filepath.ToSlash(relativeFilePath)

		newlyScannedPaths[relativeFilePath] = true

		trackNum, _ := t.Track()
		duration := 0 // Placeholder for duration, as before. You might want to get this from tag.

		newSong := models.Song{
			Title:           t.Title(),
			Artist:          t.Artist(),
			Album:           t.Album(),
			Genre:           t.Genre(),
			Year:            t.Year(),
			TrackNumber:     trackNum,
			DurationSeconds: duration,
			FilePath:        relativeFilePath,
			Filename:        filepath.Base(path),
		}

		if pic := t.Picture(); pic != nil {
			img, _, err := image.Decode(bytes.NewReader(pic.Data))
			if err != nil {
				log.Printf("Error decoding album art for %s: %v", newSong.Title, err)
				result.Errors = append(result.Errors, fmt.Sprintf("Error decoding album art for %s: %v", newSong.Title, err))
			} else {
				artTargetDir := filepath.Join(musicDir, filepath.Dir(newSong.FilePath))
				artFilePath := filepath.Join(artTargetDir, "thumb.jpg")

				// Ensure the target directory for the album art exists (e.g., Artist/Album folder)
				if err := os.MkdirAll(artTargetDir, 0755); err != nil {
					log.Printf("Error creating album art directory %s: %v", artTargetDir, err)
					result.Errors = append(result.Errors, fmt.Sprintf("Error creating album art directory %s: %v", artTargetDir, err))
				} else {
					outFile, err := os.Create(artFilePath) // os.Create truncates/overwrites if file exists
					if err != nil {
						log.Printf("Error creating album art file %s: %v", artFilePath, err)
						result.Errors = append(result.Errors, fmt.Sprintf("Error creating album art file %s: %v", artFilePath, err))
					} else {
						defer outFile.Close()
						if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90}); err != nil {
							log.Printf("Error encoding album art to JPEG for %s: %v", newSong.Title, err)
							result.Errors = append(result.Errors, fmt.Sprintf("Error encoding album art to JPEG for %s: %v", newSong.Title, err))
						}
					}
				}
				// newSong.AlbumArtPath should be the relative path from MUSIC_DIRECTORY to thumb.jpg
				newSong.AlbumArtPath = filepath.ToSlash(filepath.Join(filepath.Dir(newSong.FilePath), "thumb.jpg"))
			}
		} else {
			// If no picture embedded, ensure AlbumArtPath is empty or null in DB
			newSong.AlbumArtPath = ""
		}

		if existingSong, found := existingSongMap[newSong.FilePath]; found {
			newSong.ID = existingSong.ID
			if existingSong.Equals(&newSong) {
			} else {
				if err := DB.WithContext(ctx).Save(&newSong).Error; err != nil {
					log.Printf("Error updating song %s (ID: %s) in DB: %v", newSong.Title, newSong.ID.String(), err)
					result.Errors = append(result.Errors, fmt.Sprintf("Error updating song %s (ID: %s) in DB: %v", newSong.Title, newSong.ID.String(), err))
				} else {
					result.Updated++
				}
			}
			delete(existingSongMap, newSong.FilePath)
		} else {
			if err := DB.WithContext(ctx).Create(&newSong).Error; err != nil {
				log.Printf("Error adding new song %s to DB: %v", newSong.Title, err)
				result.Errors = append(result.Errors, fmt.Sprintf("Error adding new song %s to DB: %v", newSong.Title, err))
			} else {
				result.Added++
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error during file system walk: %w", err)
	}

	for _, songToRemove := range existingSongMap {
		if err := DB.WithContext(ctx).Delete(&songToRemove).Error; err != nil {
			log.Printf("Error deleting song %s (ID: %s) from DB: %v", songToRemove.Title, songToRemove.ID.String(), err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error deleting song %s (ID: %s) from DB: %v", songToRemove.Title, songToRemove.ID.String(), err))
		} else {
			result.Removed++
		}
	}

	return result, nil
}
