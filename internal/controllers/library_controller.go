package controllers

import (
	"context"
	"net/http"

	"github.com/demispreviotto/cajitamusical/backend/internal/db" // Replace with your module path
	"github.com/gin-gonic/gin"
)

// GetLibrary retrieves the list of songs.
func GetLibrary(c *gin.Context) {
	songs, err := db.GetSongLibrary(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song library"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}
