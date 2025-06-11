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

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/dhowden/tag"
)

type MusicScanResult struct {
	Added   int
	Updated int
	Removed int
	Errors  []string
}

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

// ScanAndStoreSongs escanea el directorio de música, lee los metadatos de las canciones y los almacena/actualiza en la base de datos.
func ScanAndStoreSongs(ctx context.Context) (*MusicScanResult, error) {
	musicDir := os.Getenv("MUSIC_DIRECTORY")
	albumArtDir := os.Getenv("ALBUM_ART_DIRECTORY")
	result := &MusicScanResult{}

	if musicDir == "" {
		return nil, fmt.Errorf("MUSIC_DIRECTORY environment variable not set")
	}
	if albumArtDir == "" {
		return nil, fmt.Errorf("ALBUM_ART_DIRECTORY environment variable not set")
	}

	// Asegurarse de que el directorio de carátulas exista
	if err := os.MkdirAll(albumArtDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create album art directory: %w", err)
	}

	// Obtener todas las canciones actuales de la DB
	var existingSongs []models.Song
	if err := DB.WithContext(ctx).Find(&existingSongs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch existing songs from DB: %w", err)
	}
	// Crear un mapa para un acceso rápido a las rutas de archivo existentes
	existingSongMap := make(map[string]models.Song)
	for _, song := range existingSongs {
		existingSongMap[song.FilePath] = song
	}

	// Escanear el sistema de archivos
	newlyScannedPaths := make(map[string]bool) // Para rastrear las rutas que encontramos en el disco

	err := filepath.Walk(musicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error accessing path %s: %v", path, err))
			return nil // No detener el walk por un error individual, solo loggear
		}
		if info.IsDir() {
			return nil
		}

		// Solo procesar archivos de audio (MP3, FLAC, M4A)
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".mp3" && ext != ".flac" && ext != ".m4a" {
			return nil
		}

		// Abrir el archivo para leer las etiquetas
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file %s: %v", path, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error opening file %s: %v", path, err))
			return nil
		}
		defer f.Close() // Asegurarse de cerrar el archivo

		// Leer las etiquetas de la canción
		t, err := tag.ReadFrom(f)
		if err != nil {
			log.Printf("Error reading tags from %s: %v", path, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error reading tags from %s: %v", path, err))
			return nil
		}

		// Construir la ruta relativa al MUSIC_DIRECTORY
		normalizedMusicDir := musicDir
		if !strings.HasSuffix(normalizedMusicDir, string(os.PathSeparator)) {
			normalizedMusicDir += string(os.PathSeparator)
		}
		relativeFilePath := strings.TrimPrefix(path, normalizedMusicDir)
		relativeFilePath = filepath.ToSlash(relativeFilePath) // Normalizar a barras inclinadas para consistencia en la DB

		newlyScannedPaths[relativeFilePath] = true // Marcar esta ruta como escaneada

		// Corrección para t.Track(): Espera dos valores de retorno (int, int).
		trackNum, _ := t.Track()

		// Corrección para duración: Si Length() no existe en tag.Metadata ni tag.Format, la inicializamos a 0 y loggeamos una advertencia.
		duration := 0

		// Declaración de hash y hashErr fuera del if para ser accesible después.
		var calculatedHash string
		var hashErr error

		newSong := models.Song{
			Title:           t.Title(),
			Artist:          t.Artist(),
			Album:           t.Album(),
			Genre:           t.Genre(),
			Year:            t.Year(),
			TrackNumber:     trackNum,
			DurationSeconds: duration, // Puede ser 0 si no se puede obtener
			FilePath:        relativeFilePath,
			Filename:        filepath.Base(path),
		}

		// Manejar la carátula del álbum
		if pic := t.Picture(); pic != nil { // 'pic' se declara aquí
			img, _, err := image.Decode(bytes.NewReader(pic.Data)) // Usa bytes.NewReader
			if err != nil {
				log.Printf("Error decoding album art for %s: %v", newSong.Title, err)
				result.Errors = append(result.Errors, fmt.Sprintf("Error decoding album art for %s: %v", newSong.Title, err))
			} else {
				// Calcula el hash aquí, donde 'pic' está disponible
				if pic.Data != nil {
					hashReader := bytes.NewReader(pic.Data)
					calculatedHash, hashErr = tag.Sum(hashReader) // tag.Sum() espera un io.Reader
				} else {
					hashErr = fmt.Errorf("picture data is nil")
				}

				hashToUse := calculatedHash
				if hashErr != nil || calculatedHash == "" {
					log.Printf("Error hashing album art data for %s: %v. Using placeholder hash.", newSong.Title, hashErr)
					hashToUse = "unknown_hash"
				}

				artFilename := fmt.Sprintf("%s.jpg", hashToUse)
				artFilePath := filepath.Join(albumArtDir, artFilename)

				// Guardar la carátula solo si no existe ya para evitar reescribir
				if _, err := os.Stat(artFilePath); os.IsNotExist(err) {
					outFile, err := os.Create(artFilePath)
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
				// Guardar la ruta relativa de la carátula para ser servida por el frontend
				newSong.AlbumArtPath = filepath.ToSlash(filepath.Join(filepath.Base(albumArtDir), artFilename))
			}
		}

		// UPSERT: Insertar o actualizar la canción en la base de datos
		if existingSong, found := existingSongMap[newSong.FilePath]; found {
			newSong.ID = existingSong.ID
			if err := DB.WithContext(ctx).Save(&newSong).Error; err != nil {
				log.Printf("Error updating song %s (ID: %d) in DB: %v", newSong.Title, newSong.ID, err)
				result.Errors = append(result.Errors, fmt.Sprintf("Error updating song %s (ID: %d) in DB: %v", newSong.Title, newSong.ID, err))
			} else {
				result.Updated++
				delete(existingSongMap, newSong.FilePath)
			}
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

	// Eliminar canciones de la DB que ya no están en el sistema de archivos
	for _, songToRemove := range existingSongMap {
		if err := DB.WithContext(ctx).Delete(&songToRemove).Error; err != nil {
			log.Printf("Error deleting song %s (ID: %d) from DB: %v", songToRemove.Title, songToRemove.ID, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Error deleting song %s (ID: %d) from DB: %v", songToRemove.Title, songToRemove.ID, err))
		} else {
			result.Removed++
		}
	}

	return result, nil
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
