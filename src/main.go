package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	engine := NewEngine()

	engine.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
}
