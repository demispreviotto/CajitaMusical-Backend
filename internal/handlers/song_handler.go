package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SongHandler is a handler for song-related operations.
type SongHandler struct {
	songService services.SongServicer
}

// NewSongHandler creates a new instance of SongHandler.
func NewSongHandler(songService services.SongServicer) *SongHandler {
	return &SongHandler{songService: songService}
}

// GetLibrary retrieves the song library.
func (h *SongHandler) GetLibrary(c *gin.Context) {
	songs, err := h.songService.GetLibrary(c.Request.Context()) // Use request context
	if err != nil {
		log.Printf("Handler: Failed to fetch song library: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

// ServeAudio serves an audio file by song ID.
func (h *SongHandler) ServeAudio(c *gin.Context) {
	songID := c.Param("songID") // Expecting song ID here

	filePath, err := h.songService.GetSongFilePath(songID)
	if err != nil {
		log.Printf("Handler: Error serving audio for song ID %s: %v", songID, err)
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "path traversal") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"}) // Prevent leaking internal paths
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audio file"})
		}
		return
	}

	c.File(filePath)
}

// ServeAlbumArt serves an album art image.
func (h *SongHandler) ServeAlbumArt(c *gin.Context) {
	imageFileName := c.Param("filename") // This should be the hash.jpg filename

	imagePath, err := h.songService.GetAlbumArtPath(imageFileName)
	if err != nil {
		log.Printf("Handler: Error serving album art %s: %v", imageFileName, err)
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "path traversal") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image path"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album art"})
		}
		return
	}

	c.File(imagePath)
}

// ScanMusicLibrary triggers a manual scan of the music directory.
func (h *SongHandler) ScanMusicLibrary(c *gin.Context) {
	result, err := h.songService.ScanMusicLibrary(c.Request.Context())
	if err != nil {
		log.Printf("Handler: Failed to trigger music scan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Music library scan initiated successfully",
		"result":  result,
	})
}
