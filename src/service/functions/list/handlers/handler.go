package handlers

import (
	"io/ioutil"
	"net/http"
	"service/common"
	"service/common/graph"

	"encoding/xml"

	"github.com/gin-gonic/gin"
)

type route struct {
	Flights []*flightItem `json:"flights"`
}

type flightItem struct {
	Flight  *common.Flight  `json:"flight"`
	Pricing *common.Pricing `json:"pricing"`
}

func (f *flightItem) IsAccessibleFrom(from interface{}) bool {
	return (from.(*flightItem)).Flight.ArrivalTimeStamp.Before(f.Flight.DepartureTimeStamp.Time) //TBD: add some time reserve?
}

func (f *flightItem) Cost() float32 {
	return 0
}

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

	var flights []flightItem
	nodes := make(map[string]bool)

	for p, f := range data.PricedItineraries.Flights {
		for _, items := range []common.PricedItinerary{f.OnwardPricedItinerary, f.ReturnPricedItinerary} {

			for idx, flight := range items.Flights.Flight {
				flights = append(flights, flightItem{
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

	paths := g.GetPaths(req.Source, req.Destination, req.MaxFlightsInRoute)
	var routes []route

	for _, p := range paths {
		var flights []*flightItem
		edges := p.Edges()
		for _, edge := range edges {
			flights = append(flights, edge.(*flightItem))
		}
		routes = append(routes, route{Flights: flights})
	}

	/*fmt.Printf("Got %d paths", len(paths))
	fmt.Println()
	for idx, r := range paths {
		if idx >= 0 {
			edges := r.Edges()
			fmt.Printf("Got %d edges in route %d", len(edges), idx)
			fmt.Println()
			for _, edge := range edges {
				edge := edge.(*flightItem)
				fl := edge.Flight
				fmt.Printf("%s [%s]: %s (%s) ==> %s (%s)", fl.Carrier.Name, fl.FlightNumber, fl.Source, fl.DepartureTimeStamp, fl.Destination, fl.ArrivalTimeStamp)
				fmt.Println()
			}
			fmt.Println()
		}
	}*/

	c.JSON(http.StatusOK, gin.H{"success": true, "routes": routes})
}
