package server

import (
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/gin-gonic/gin"
)

//Start running server based on env PLATFORM
func Start(engine *gin.Engine) {
	address := ":3000"
	platform := os.Getenv("PLATFORM")

	if platform == "aws_lambda" {
		log.Fatal(gateway.ListenAndServe(address, engine))
	} else {
		log.Fatal(http.ListenAndServe(address, engine))
	}
}
