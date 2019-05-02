package common

import (
	"encoding/xml"
	"mime/multipart"
	"time"
)

type SingleDataRequest struct {
	Data              *multipart.FileHeader `form:"data" binding:"required"`
	Source            string                `form:"source" binding:"required"`
	Destination       string                `form:"destination" binding:"required"`
	MaxFlightsInRoute int                   `form:"max_flights_in_route"`
}

type Timestamp struct {
	time.Time
}

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

type Flight struct {
	Carrier struct {
		Name string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"Carrier"`
	FlightNumber       string    `xml:"FlightNumber"`
	Source             string    `xml:"Source"`
	Destination        string    `xml:"Destination"`
	DepartureTimeStamp Timestamp `xml:"DepartureTimeStamp"`
	ArrivalTimeStamp   Timestamp `xml:"ArrivalTimeStamp"`
	Class              string    `xml:"Class"`
	NumberOfStops      int       `xml:"NumberOfStops"`
	FareBasis          string    `xml:"FareBasis"`
	WarningText        string    `xml:"WarningText"`
	TicketType         string    `xml:"TicketType"`
}

type Flights struct {
	Flight []Flight `xml:"Flight"`
}

type PricedItinerary struct {
	Flights Flights `xml:"Flights"`
}

type Pricing struct {
	Currency       string `xml:"currency,attr"`
	ServiceCharges []struct {
		Amount     float32 `xml:",chardata"`
		Type       string  `xml:"type,attr"`
		ChargeType string  `xml:"ChargeType,attr"`
	} `xml:"ServiceCharges"`
}

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
