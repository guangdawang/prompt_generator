package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func respondInternalError(c *gin.Context, err error) {
	log.Printf("internal error: %v", err)
	respondError(c, http.StatusInternalServerError, "internal server error")
}
