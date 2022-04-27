package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	authzService := NewAuthzService()

	go (func(a *AuthzService) {
		for {
			a.Refresh()
			time.Sleep(20 * time.Minute)
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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
}
