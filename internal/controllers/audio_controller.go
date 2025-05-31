package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ServeAudio(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	musicDir := os.Getenv("MUSIC_DIRECTORY")
	if musicDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MUSIC_DIRECTORY not configured"})
		return
	}

	filePath := filepath.Join(musicDir, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio file not found"})
		return
	}

	// Serve the content
	c.File(filePath)
}
