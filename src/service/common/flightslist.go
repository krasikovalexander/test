package common

import (
	"log"
	"sync"
	"time"

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
	start := time.Now()
	jobs := make(chan comparePair)
	updates := make(chan FlightUpdate)

	wgResults := sync.WaitGroup{}
	wgResults.Add(1)
	go func() {
		defer wgResults.Done()
		for result := range updates {
			modifications = append(modifications, result)
		}
	}()

	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(WorkersCount)

	var worker = func(jobs <-chan comparePair, results chan<- FlightUpdate) {
		defer wgWorkers.Done()
		for compare := range jobs {
			changelog, _ := diff.Diff(compare.flightA, compare.flightB)
			results <- FlightUpdate{compare.flightA, changelog}
		}
	}

	for w := 0; w < WorkersCount; w++ {
		go worker(jobs, updates)
	}

	for _, flightA := range flightsA.flightItems {
		if flightB, exist := flightsB.flightItems[flightA.Flight.Key()]; exist {
			if *flightA.Flight != *flightB.Flight {
				jobs <- comparePair{flightA, flightB}
			}
		} else {
			removals = append(removals, flightA)
		}
	}
	close(jobs)
	wgWorkers.Wait()
	close(updates)
	wgResults.Wait()

	log.Printf("WorkersCount %d took %s", WorkersCount, time.Since(start))

	for _, flightB := range flightsB.flightItems {
		if _, exist := flightsA.flightItems[flightB.Flight.Key()]; !exist {
			additions = append(additions, flightB)
		}
	}
	return
}
