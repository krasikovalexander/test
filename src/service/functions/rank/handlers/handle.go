package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rank handler"})
}
