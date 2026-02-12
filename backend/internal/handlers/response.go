package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, status int, message string) {
	// 对于 4xx 错误，返回更详细的错误信息以便前端展示。
	// 对于 5xx 错误，返回通用错误信息，避免泄露内部实现细节。
	if status >= 500 {
		c.JSON(status, gin.H{"error": http.StatusText(status)})
	} else {
		c.JSON(status, gin.H{"error": http.StatusText(status), "message": message})
	}
}

func respondInternalError(c *gin.Context, err error) {
	log.Printf("internal error: %v", err)
	respondError(c, http.StatusInternalServerError, "internal server error")
}
