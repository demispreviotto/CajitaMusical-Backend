package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// songHandler is a handler for song-related operations.
type songHandler struct {
	songService services.SongServicer
}

// NewsongHandler creates a new instance of songHandler.
func NewsongHandler(songService services.SongServicer) *songHandler {
	return &songHandler{songService: songService}
}

// GetLibrary retrieves the song library.
func (h *songHandler) GetLibrary(c *gin.Context) {
	songs, err := h.songService.GetLibrary(c.Request.Context()) // Use request context
	if err != nil {
		log.Printf("Handler: Failed to fetch song library: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

// ServeAudio serves an audio file by song ID.
func (h *songHandler) ServeAudio(c *gin.Context) {
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
func (h *songHandler) ServeAlbumArt(c *gin.Context) {
	// The filename parameter will now be the full relative path, e.g., "Artist/Album/thumb.jpg"
	imageRelativePath := c.Param("filepath") // Use "filepath" as the URL parameter name for clarity
	if imageRelativePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album art path is missing"})
		return
	}

	// Decode URL-encoded parts, if any
	decodedPath := filepath.FromSlash(imageRelativePath)

	musicDir := os.Getenv("MUSIC_DIRECTORY")
	if musicDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MUSIC_DIRECTORY not set"})
		return
	}

	fullExpectedPath := filepath.Join(musicDir, decodedPath)
	absMusicDir, err := filepath.Abs(musicDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve music directory"})
		return
	}
	absRequestedPath, err := filepath.Abs(fullExpectedPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path requested"})
		return
	}

	if !strings.HasPrefix(absRequestedPath, absMusicDir) {
		log.Printf("Attempted path traversal in album art request: %s", imageRelativePath)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if _, err := os.Stat(fullExpectedPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album art not found"})
		return
	} else if err != nil {
		log.Printf("Error accessing album art file %s: %v", fullExpectedPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read album art"})
		return
	}

	c.File(fullExpectedPath)
}

// ScanMusicLibrary triggers a manual scan of the music directory.
func (h *songHandler) ScanMusicLibrary(c *gin.Context) {
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
