package handlers

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"service/common"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	var req common.CompareDataRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	file, _ := req.DataA.Open()
	content, _ := ioutil.ReadAll(file)

	var dataA common.AirFareSearchResponse
	if err := xml.Unmarshal(content, &dataA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	file, _ = req.DataB.Open()
	content, _ = ioutil.ReadAll(file)

	var dataB common.AirFareSearchResponse
	if err := xml.Unmarshal(content, &dataB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	flightsA := common.NewFlightsList(&dataA)
	flightsB := common.NewFlightsList(&dataB)

	response := gin.H{"success": true}

	response["additions"],
		response["removals"],
		response["updates"] = flightsB.Diff(flightsA)

	c.JSON(http.StatusOK, response)
}
