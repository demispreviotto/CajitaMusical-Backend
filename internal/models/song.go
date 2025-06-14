package models

import (
	"time"

	"github.com/google/uuid"
)

// Song representa una canción almacenada en la base de datos.
type Song struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"` // ID único de la canción en la DB, ahora UUID
	Title           string    `gorm:"size:255;not null" json:"title"`
	Artist          string    `gorm:"size:255;not null" json:"artist"`
	Album           string    `gorm:"size:255" json:"album,omitempty"`
	TrackNumber     int       `json:"track_number,omitempty"`
	Genre           string    `gorm:"size:255" json:"genre,omitempty"`
	Year            int       `json:"year,omitempty"`
	DurationSeconds int       `json:"duration_seconds"`                          // Duración en segundos
	FilePath        string    `gorm:"size:512;unique;not null" json:"file_path"` // FilePath es la ruta del archivo relativa al MUSIC_DIRECTORY. Ej: "Artista/Álbum/Cancion.mp3". Debe ser único para evitar duplicados del mismo archivo.
	Filename        string    `gorm:"size:255;not null" json:"filename"`         // ej: "Cancion.mp3"
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	AlbumArtPath    string    `gorm:"size:512" json:"album_art_path,omitempty"`
}

func (s *Song) Equals(other *Song) bool {
	return s.Title == other.Title &&
		s.Artist == other.Artist &&
		s.Album == other.Album &&
		s.TrackNumber == other.TrackNumber &&
		s.Genre == other.Genre &&
		s.Year == other.Year &&
		s.DurationSeconds == other.DurationSeconds &&
		s.AlbumArtPath == other.AlbumArtPath
	// You might also compare other fields if you add them (e.g., Copyright, Composer)
}
