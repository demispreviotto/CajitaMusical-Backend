package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/services" // Import the services package
	"github.com/gin-gonic/gin"
)

// SongHandler is a handler for song-related operations.
type SongHandler struct {
	songService services.SongServicer // Dependency on the song service
}

// NewSongHandler creates a new instance of SongHandler.
func NewSongHandler(songService services.SongServicer) *SongHandler {
	return &SongHandler{songService: songService}
}

// GetLibrary retrieves the song library.
func (h *SongHandler) GetLibrary(c *gin.Context) {
	songs, err := h.songService.GetLibrary(context.Background())
	if err != nil {
		log.Printf("Handler: Failed to fetch song library: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return the service error
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

// ServeAudio serves an audio file from the MUSIC_DIRECTORY.
func (h *SongHandler) ServeAudio(c *gin.Context) {
	filename := c.Param("filename")

	// Use the service to get the validated file path
	filePath, err := h.songService.GetSongFilePath(filename)
	if err != nil {
		log.Printf("Handler: Error serving audio file %s: %v", filename, err)
		// Map service errors to appropriate HTTP responses
		if err.Error() == "filename is required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "audio file not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audio file"})
		}
		return
	}

	// Serve the content directly from the handler, as Gin's c.File is HTTP-specific
	c.File(filePath)
}
