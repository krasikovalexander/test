package common

import (
	"encoding/xml"
	"mime/multipart"
	"service/common/graph"
	"time"
)

//SingleDataRequest is a multipart/form-data binding
type SingleDataRequest struct {
	Data              *multipart.FileHeader `form:"data" binding:"required"`
	Source            string                `form:"source" binding:"required"`
	Destination       string                `form:"destination" binding:"required"`
	MaxFlightsInRoute int                   `form:"max_flights_in_route"`
}

//Timestamp time.Time with unmarshal 2006-01-02T1504 support
type Timestamp struct {
	time.Time
}

//UnmarshalXML "2006-01-02T1504" to Timestamp
func (t *Timestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const format = "2006-01-02T1504"
	var str string
	d.DecodeElement(&str, &start)
	parsed, err := time.Parse(format, str)
	if err != nil {
		return err
	}
	*t = Timestamp{parsed}
	return nil
}

//Flight information
type Flight struct {
	Carrier struct {
		Name string `xml:",chardata" json:"name"`
		ID   string `xml:"id,attr"  json:"id"`
	} `xml:"Carrier"  json:"carrier"`
	FlightNumber       string    `xml:"FlightNumber"  json:"flightNumber"`
	Source             string    `xml:"Source"  json:"source"`
	Destination        string    `xml:"Destination"  json:"destination"`
	DepartureTimeStamp Timestamp `xml:"DepartureTimeStamp"  json:"departureTimeStamp"`
	ArrivalTimeStamp   Timestamp `xml:"ArrivalTimeStamp"  json:"arrivalTimeStamp"`
	Class              string    `xml:"Class"  json:"class"`
	NumberOfStops      int       `xml:"NumberOfStops"  json:"numberOfStops"`
	FareBasis          string    `xml:"FareBasis"  json:"fareBasis"`
	WarningText        string    `xml:"WarningText"  json:"warningText"`
	TicketType         string    `xml:"TicketType"  json:"ticketType"`
}

//Flights accessory structure
type Flights struct {
	Flight []Flight `xml:"Flight"`
}

//PricedItinerary accessory structure
type PricedItinerary struct {
	Flights Flights `xml:"Flights"`
}

//Pricing accessory structure
type Pricing struct {
	Currency       string `xml:"currency,attr"  json:"currency"`
	ServiceCharges []struct {
		Amount     float32 `xml:",chardata"  json:"amount"`
		Type       string  `xml:"type,attr"  json:"type"`
		ChargeType string  `xml:"ChargeType,attr"  json:"chargeType"`
	} `xml:"ServiceCharges"  json:"serviceCharges"`
}

//GetTotalAmount returns total flight cost
func (p *Pricing) GetTotalAmount() (amount float32, ok bool) {
	if p == nil {
		return
	}
	for _, charge := range p.ServiceCharges {
		if charge.ChargeType == "TotalAmount" { //TBD: currency issues
			return charge.Amount, true
		}
	}
	return
}

//AirFareSearchResponse xml binding
type AirFareSearchResponse struct {
	RequestTime       string `xml:"RequestTime,attr"`
	ResponseTime      string `xml:"ResponseTime,attr"`
	RequestID         string `xml:"RequestId"`
	PricedItineraries struct {
		Flights []struct {
			OnwardPricedItinerary PricedItinerary `xml:"OnwardPricedItinerary"`
			ReturnPricedItinerary PricedItinerary `xml:"ReturnPricedItinerary"`

			Pricing Pricing `xml:"Pricing"`
		} `xml:"Flights"`
	} `xml:"PricedItineraries"`
}

//TransferTimeInMinutes time window between transshipping
const TransferTimeInMinutes = 60

//Route is a list of FlightItem
type Route struct {
	Flights []*FlightItem `json:"flights"`
}

//FlightItem stores info about Flight and it's Pricing
type FlightItem struct {
	Flight  *Flight  `json:"flight"`
	Pricing *Pricing `json:"pricing"`
}

//IsAccessibleFrom detects if flight is available due arrival and departure time
func (f *FlightItem) IsAccessibleFrom(from interface{}) bool {
	nextFlightDepartureTime := f.Flight.DepartureTimeStamp.Time.Add(-TransferTimeInMinutes * time.Minute) //need some time for transshipment
	return (from.(*FlightItem)).Flight.ArrivalTimeStamp.Before(nextFlightDepartureTime)
}

//NewFlightsGraph creates graph by data from AirFareSearchResponse
func NewFlightsGraph(data *AirFareSearchResponse) *graph.Graph {
	var flights []FlightItem
	nodes := make(map[string]bool)

	for p, f := range data.PricedItineraries.Flights {
		for _, items := range []PricedItinerary{f.OnwardPricedItinerary, f.ReturnPricedItinerary} {

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
	return g
}
