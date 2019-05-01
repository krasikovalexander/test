package handlers

import (
	"io/ioutil"
	"net/http"
	"service/common/server"

	"encoding/xml"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	var req server.SingleDataRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	file, _ := req.Data.Open()
	content, _ := ioutil.ReadAll(file)

	var data server.AirFareSearchResponse
	if err := xml.Unmarshal(content, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "list handler"})
}
