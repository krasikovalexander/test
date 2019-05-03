package handlers

import (
	"service/common"
	"service/common/graph"
	"time"
)

type (
	Criterion struct {
		Paths    []*graph.Path
		hasValue bool
		Value    interface{}
		Fn       func(c *Criterion, path *graph.Path) (interface{}, bool)
	}
)

func (c *Criterion) GetResult() []*graph.Path {
	return c.Paths
}

func (c *Criterion) Apply(path *graph.Path) {
	value, isOptimal := c.Fn(c, path)

	if isOptimal {
		if c.hasValue && c.Value != value {
			c.Paths = nil
		}
		c.Paths = append(c.Paths, path)
		c.Value = value
		c.hasValue = true
	}
}

func NewMinimumCostCriterion() *Criterion {
	return &Criterion{
		Fn: func(c *Criterion, path *graph.Path) (interface{}, bool) {
			var totalCost float32

			for _, item := range path.Edges() {
				item := item.(*common.FlightItem)
				if price, ok := item.Pricing.GetTotalAmount(); ok {
					totalCost = totalCost + price
					if c.hasValue && totalCost > c.Value.(float32) {
						return 0, false
					}
				}
			}
			return totalCost, true
		},
	}
}

func NewMaximumCostCriterion() *Criterion {
	return &Criterion{
		Fn: func(c *Criterion, path *graph.Path) (interface{}, bool) {
			var totalCost float32

			for _, item := range path.Edges() {
				item := item.(*common.FlightItem)
				if price, ok := item.Pricing.GetTotalAmount(); ok {
					totalCost = totalCost + price
				}
			}
			if c.hasValue && totalCost < c.Value.(float32) {
				return 0, false
			}
			return totalCost, true
		},
	}
}

func NewMinimumTimeCriterion() *Criterion {
	return &Criterion{
		Fn: func(c *Criterion, path *graph.Path) (interface{}, bool) {
			edges := path.Edges()

			departureTime := edges[0].(*common.FlightItem).Flight.DepartureTimeStamp.Time
			arrivalTime := edges[len(edges)-1].(*common.FlightItem).Flight.ArrivalTimeStamp.Time

			totalTime := arrivalTime.Sub(departureTime)

			if c.hasValue && totalTime > c.Value.(time.Duration) {
				return 0, false
			}
			return totalTime, true
		},
	}
}

func NewMaximumTimeCriterion() *Criterion {
	return &Criterion{
		Fn: func(c *Criterion, path *graph.Path) (interface{}, bool) {
			edges := path.Edges()

			departureTime := edges[0].(*common.FlightItem).Flight.DepartureTimeStamp.Time
			arrivalTime := edges[len(edges)-1].(*common.FlightItem).Flight.ArrivalTimeStamp.Time

			totalTime := arrivalTime.Sub(departureTime)

			if c.hasValue && totalTime < c.Value.(time.Duration) {
				return 0, false
			}
			return totalTime, true
		},
	}
}
