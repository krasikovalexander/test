package main

import (
	"service/common/server"
	"service/functions/compare/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.POST("/compare", handlers.Handle)

	server.Start(router)
}
