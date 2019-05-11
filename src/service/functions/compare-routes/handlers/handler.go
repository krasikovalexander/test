package handlers

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"service/common"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	var req common.CompareRoutesDataRequest

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

	graphA := common.NewFlightsGraph(&dataA)
	graphB := common.NewFlightsGraph(&dataB)

	pathsA := graphA.GetPaths(req.Source, req.Destination, req.MaxFlightsInRoute)
	pathsB := graphB.GetPaths(req.Source, req.Destination, req.MaxFlightsInRoute)

	var routesA []common.Route
	var routesB []common.Route

	for _, p := range pathsA {
		var flights []*common.FlightItem
		edges := p.Edges()
		for _, edge := range edges {
			flights = append(flights, edge.(*common.FlightItem))
		}
		routesA = append(routesA, common.Route{Flights: flights})
	}

	for _, p := range pathsB {
		var flights []*common.FlightItem
		edges := p.Edges()
		for _, edge := range edges {
			flights = append(flights, edge.(*common.FlightItem))
		}
		routesB = append(routesB, common.Route{Flights: flights})
	}

	listA := common.NewRoutesList(routesA)
	listB := common.NewRoutesList(routesB)

	response := gin.H{"success": true}

	response["additions"],
		response["removals"],
		response["updates"] = listB.Diff(listA)

	c.JSON(http.StatusOK, response)
}
