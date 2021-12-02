package main

import (
	"io/ioutil"
	"log"
	"net/http"
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

	// don't sleep
	go (func() {
		for {
			resp, _ := http.Get(app)
			body, _ := ioutil.ReadAll(resp.Body)
			log.Println(string(body))
			time.Sleep(25 * time.Minute)
		}
	})()

	authzService.r.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
}
