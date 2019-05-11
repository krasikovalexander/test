package common

import (
	"encoding/xml"
	"fmt"
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

//CompareDataRequest is a multipart/form-data binding
type CompareDataRequest struct {
	DataA *multipart.FileHeader `form:"data_a" binding:"required"`
	DataB *multipart.FileHeader `form:"data_b" binding:"required"`
}

//CompareRoutesDataRequest is a multipart/form-data binding
type CompareRoutesDataRequest struct {
	DataA *multipart.FileHeader `form:"data_a" binding:"required"`
	DataB *multipart.FileHeader `form:"data_b" binding:"required"`

	Source            string `form:"source" binding:"required"`
	Destination       string `form:"destination" binding:"required"`
	MaxFlightsInRoute int    `form:"max_flights_in_route"`
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
		Name string `xml:",chardata" json:"name" diff:"name"`
		ID   string `xml:"id,attr" json:"id" diff:"id"`
	} `xml:"Carrier" json:"carrier" diff:"carrier"`
	FlightNumber       string    `xml:"FlightNumber" json:"flightNumber" diff:"flightNumber"`
	Source             string    `xml:"Source" json:"source" diff:"source"`
	Destination        string    `xml:"Destination" json:"destination" diff:"destination"`
	DepartureTimeStamp Timestamp `xml:"DepartureTimeStamp" json:"departureTimeStamp" diff:"departureTimeStamp"`
	ArrivalTimeStamp   Timestamp `xml:"ArrivalTimeStamp" json:"arrivalTimeStamp" diff:"arrivalTimeStamp"`
	Class              string    `xml:"Class" json:"class" diff:"class"`
	NumberOfStops      int       `xml:"NumberOfStops" json:"numberOfStops" diff:"numberOfStops"`
	FareBasis          string    `xml:"FareBasis" json:"fareBasis" diff:"fareBasis"`
	WarningText        string    `xml:"WarningText" json:"warningText" diff:"warningText"`
	TicketType         string    `xml:"TicketType" json:"ticketType"  diff:"ticketType"`
}

func (f *Flight) Key() string {
	return fmt.Sprintf("%s:%s:%s:%s", f.Carrier.Name, f.FlightNumber, f.DepartureTimeStamp.Format("01-02-2006"), f.FareBasis)
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
	Currency       string `xml:"currency,attr"  json:"currency" diff:"currency"`
	ServiceCharges []struct {
		Amount     float32 `xml:",chardata"  json:"amount" diff:"amount"`
		Type       string  `xml:"type,attr"  json:"type" diff:"type"`
		ChargeType string  `xml:"ChargeType,attr"  json:"chargeType" diff:"chargeType"`
	} `xml:"ServiceCharges"  json:"serviceCharges" diff:"serviceCharges"`
}

//GetTotalAmount returns total flight cost
func (p *Pricing) GetTotalAmount() (amount float32, ok bool) {
	if p == nil {
		return
	}
	for _, charge := range p.ServiceCharges {
		if charge.ChargeType == "TotalAmount" && charge.Type == "SingleAdult" { //TBD: currency issues
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
	Flights []*FlightItem `json:"flights" diff:"flights"`
}

//Key returns Route's Flights composite key
func (r *Route) Key() (key string) {
	for _, f := range r.Flights {
		key = fmt.Sprintf("%s:%s", key, f.Flight.Key())
	}
	return
}

//FlightItem stores info about Flight and it's Pricing
type FlightItem struct {
	Flight  *Flight  `json:"flight" diff:"flight"`
	Pricing *Pricing `json:"pricing" diff:"pricing"`
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
