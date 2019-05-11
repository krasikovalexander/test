package common

import (
	"sync"

	"github.com/r3labs/diff"
)

var WorkersCount int = 10

type FlightsList struct {
	flightItems map[string]FlightItem
}

//NewFlightsList creates FlightsList by data from AirFareSearchResponse
func NewFlightsList(data *AirFareSearchResponse) *FlightsList {
	fl := FlightsList{
		flightItems: make(map[string]FlightItem),
	}
	for p, f := range data.PricedItineraries.Flights {
		for _, items := range []PricedItinerary{f.OnwardPricedItinerary, f.ReturnPricedItinerary} {

			for idx, flight := range items.Flights.Flight {
				fl.flightItems[flight.Key()] = FlightItem{
					Flight:  &items.Flights.Flight[idx],
					Pricing: &data.PricedItineraries.Flights[p].Pricing,
				}
			}
		}
	}

	return &fl
}

type FlightUpdate struct {
	FlightItem FlightItem     `json:"origin"`
	Changelog  diff.Changelog `json:"changes"`
}

func (flightsB *FlightsList) Diff(flightsA *FlightsList) (additions []FlightItem, removals []FlightItem, modifications []FlightUpdate) {

	type comparePair struct {
		flightA FlightItem
		flightB FlightItem
	}

	jobs := make(chan comparePair)
	updates := make(chan FlightUpdate)

	wgUpdates := sync.WaitGroup{}
	wgUpdates.Add(1)

	go func(updates <-chan FlightUpdate) {
		defer wgUpdates.Done()
		for result := range updates {
			modifications = append(modifications, result)
		}
	}(updates)

	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(WorkersCount)

	var worker = func(jobs <-chan comparePair, updates chan<- FlightUpdate) {
		defer wgWorkers.Done()
		for compare := range jobs {
			changelog, _ := diff.Diff(compare.flightA, compare.flightB)
			updates <- FlightUpdate{compare.flightA, changelog}
		}
	}

	for w := 0; w < WorkersCount; w++ {
		go worker(jobs, updates)
	}

	for _, flightA := range flightsA.flightItems {
		if flightB, exist := flightsB.flightItems[flightA.Flight.Key()]; exist {
			priceA, _ := flightA.Pricing.GetTotalAmount()
			priceB, _ := flightB.Pricing.GetTotalAmount()
			if *flightA.Flight != *flightB.Flight || priceA != priceB {
				jobs <- comparePair{flightA, flightB}
			}
		} else {
			removals = append(removals, flightA)
		}
	}
	close(jobs)
	wgWorkers.Wait()
	close(updates)
	wgUpdates.Wait()

	for _, flightB := range flightsB.flightItems {
		if _, exist := flightsA.flightItems[flightB.Flight.Key()]; !exist {
			additions = append(additions, flightB)
		}
	}
	return
}
