package common

import (
	"sync"

	"github.com/r3labs/diff"
)

var RoutesWorkersCount int = 10

type RoutesList struct {
	routes map[string]*Route
}

//NewRoutesList creates RoutesList by []Route
func NewRoutesList(routes []Route) *RoutesList {
	rl := RoutesList{
		routes: make(map[string]*Route),
	}

	for idx, route := range routes {
		rl.routes[route.Key()] = &routes[idx]
	}

	return &rl
}

type RouteUpdate struct {
	Route     Route          `json:"origin"`
	Changelog diff.Changelog `json:"changes"`
}

func (routesB *RoutesList) Diff(routesA *RoutesList) (additions []Route, removals []Route, modifications []RouteUpdate) {

	type comparePair struct {
		routeA Route
		routeB Route
	}

	jobs := make(chan comparePair)
	updates := make(chan RouteUpdate)

	wgUpdates := sync.WaitGroup{}
	wgUpdates.Add(1)

	go func(updates <-chan RouteUpdate) {
		defer wgUpdates.Done()
		for result := range updates {
			modifications = append(modifications, result)
		}
	}(updates)

	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(WorkersCount)

	var worker = func(jobs <-chan comparePair, updates chan<- RouteUpdate) {
		defer wgWorkers.Done()
		for compare := range jobs {
			changelog, _ := diff.Diff(compare.routeA, compare.routeB)
			if len(changelog) > 0 {
				updates <- RouteUpdate{compare.routeA, changelog}
			}
		}
	}

	for w := 0; w < WorkersCount; w++ {
		go worker(jobs, updates)
	}

	for _, routeA := range routesA.routes {
		if routeB, exist := routesB.routes[routeA.Key()]; exist {
			jobs <- comparePair{*routeA, *routeB}
		} else {
			removals = append(removals, *routeA)
		}
	}
	close(jobs)
	wgWorkers.Wait()
	close(updates)
	wgUpdates.Wait()

	for _, routeB := range routesB.routes {
		if _, exist := routesA.routes[routeB.Key()]; !exist {
			additions = append(additions, *routeB)
		}
	}
	return
}
