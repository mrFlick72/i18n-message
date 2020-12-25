package main

import (
	"github/mrflick72/i18n-message/configuration"
	"sync"
)

func main() {
	// Creates an iris application without any middleware by default
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go configuration.DocumentUpdatesListener(wg)
	go configuration.NewApplicationServer(wg)
	go configuration.NewActuatorServer(wg)

	wg.Wait()
}
