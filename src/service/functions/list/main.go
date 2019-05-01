package main

import (
	"service/common/server"
	"service/functions/list/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.POST("/list", handlers.Handle)

	server.Start(router)
}
