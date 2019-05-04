package main

import (
	"github.com/gin-gonic/gin"

	"service/common/server"
	compare "service/functions/compare/handlers"
	list "service/functions/list/handlers"
	rank "service/functions/rank/handlers"
)

func main() {
	router := gin.New()
	router.POST("/compare", compare.Handle)
	router.POST("/list", list.Handle)
	router.POST("/rank", rank.Handle)

	server.Start(router)
}
