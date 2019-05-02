package common

import (
	"encoding/xml"
	"mime/multipart"
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
