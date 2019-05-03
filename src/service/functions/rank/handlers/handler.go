package handlers

import (
	"io/ioutil"
	"net/http"
	"service/common"

	"encoding/xml"

	"github.com/gin-gonic/gin"
)

//Handle api call handler. Consumes multipart/form-data, produces json
func Handle(c *gin.Context) {
	var req common.SingleDataRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	file, _ := req.Data.Open()
	content, _ := ioutil.ReadAll(file)

	var data common.AirFareSearchResponse
	if err := xml.Unmarshal(content, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	g := common.NewFlightsGraph(&data)

	minCost := NewMinimumCostCriterion()
	maxCost := NewMaximumCostCriterion()
	minTime := NewMinimumTimeCriterion()
	maxTime := NewMaximumTimeCriterion()

	g.SearchOptimalPaths(req.Source, req.Destination, req.MaxFlightsInRoute, minCost, maxCost, minTime, maxTime)

	items := map[string]*Criterion{
		"minCost": minCost,
		"maxCost": maxCost,
		"minTime": minTime,
		"maxTime": maxTime,
	}

	result := make(map[string]interface{})

	for key, criterion := range items {
		var routes []common.Route
		paths := criterion.GetResult()
		for _, p := range paths {
			var flights []*common.FlightItem
			edges := p.Edges()
			for _, edge := range edges {
				flights = append(flights, edge.(*common.FlightItem))
			}
			routes = append(routes, common.Route{Flights: flights})
		}
		result[key] = routes
	}

	result["success"] = true
	c.JSON(http.StatusOK, result)

	/*paths := maxTime.GetResult()
	fmt.Println(maxTime.Value)
	fmt.Printf("Got %d paths", len(paths))
	fmt.Println()
	for idx, r := range paths {
		if idx >= 0 {
			edges := r.Edges()
			fmt.Printf("Got %d edges in route %d", len(edges), idx)
			fmt.Println()
			for _, edge := range edges {
				edge := edge.(*common.FlightItem)
				fl := edge.Flight
				fmt.Printf("%s [%s]: %s (%s) ==> %s (%s)", fl.Carrier.Name, fl.FlightNumber, fl.Source, fl.DepartureTimeStamp, fl.Destination, fl.ArrivalTimeStamp)
				fmt.Println()
			}
			fmt.Println()
		}
	}*/
}
