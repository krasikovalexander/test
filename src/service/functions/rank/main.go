package main

import (
	"service/common/server"
	"service/functions/rank/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.POST("/rank", handlers.Handle)

	server.Start(router)
}
