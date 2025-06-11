package handlers

import (
	"context"
	"net/http"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func GetLibrary(c *gin.Context) {
	songs, err := db.GetSongLibrary(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song library"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}
