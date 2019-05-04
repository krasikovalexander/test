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
	paths := g.GetPaths(req.Source, req.Destination, req.MaxFlightsInRoute)

	var routes []common.Route

	for _, p := range paths {
		var flights []*common.FlightItem
		edges := p.Edges()
		for _, edge := range edges {
			flights = append(flights, edge.(*common.FlightItem))
		}
		routes = append(routes, common.Route{Flights: flights})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "routes": routes})

	/*fmt.Printf("Got %d paths", len(paths))
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
