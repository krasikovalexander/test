package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"service/common"
	"service/common/graph"

	"encoding/xml"

	"github.com/gin-gonic/gin"
)

type FlightItem struct {
	Flight  *common.Flight
	Pricing *common.Pricing
}

func (f *FlightItem) IsAccessibleFrom(from interface{}) bool {
	return (from.(*FlightItem)).Flight.ArrivalTimeStamp.Before(f.Flight.DepartureTimeStamp.Time) //TBD: add some threshod?
}

func (f *FlightItem) Cost() float32 {
	return 0
}

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

	var flights []FlightItem
	nodes := make(map[string]bool)

	for p, f := range data.PricedItineraries.Flights {
		for _, items := range []common.PricedItinerary{f.OnwardPricedItinerary, f.ReturnPricedItinerary} {

			for idx, flight := range items.Flights.Flight {
				flights = append(flights, FlightItem{
					Flight:  &items.Flights.Flight[idx],
					Pricing: &data.PricedItineraries.Flights[p].Pricing,
				})
				if !nodes[flight.Source] {
					nodes[flight.Source] = true
				}
				if !nodes[flight.Destination] {
					nodes[flight.Destination] = true
				}
			}
		}
	}

	g := graph.NewGraph(len(nodes))
	for idx, item := range flights {
		g.AddEdge(item.Flight.Source, item.Flight.Destination, &flights[idx])
	}

	routes := g.GetPaths(req.Source, req.Destination, req.MaxFlightsInRoute)
	fmt.Printf("Got %d routes", len(routes))
	fmt.Println()
	for idx, r := range routes {
		if idx >= 0 {
			edges := r.Edges()
			fmt.Printf("Got %d edges in route %d", len(edges), idx)
			fmt.Println()
			for _, edge := range edges {
				edge := edge.(*FlightItem)
				fl := edge.Flight
				fmt.Printf("%s: %s (%s) ==> %s (%s)", fl.Carrier.Name, fl.Source, fl.DepartureTimeStamp, fl.Destination, fl.ArrivalTimeStamp)
				fmt.Println()
			}
			fmt.Println()
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "list handler"})
}
