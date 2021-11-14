package main

import (
	"time"
)

func main() {
	authzService := NewAuthzService()

	go (func(a *AuthzService) {
		for {
			a.Refresh()
			time.Sleep(50 * time.Minute)
		}
	})(authzService)

	go (func(a *AuthzService) {
		graphService := NewGraphService(a)

		for {
			graphService.FetchEvents()
			time.Sleep(checkEventsDuration)
		}
	})(authzService)

	authzService.r.Run()
}
