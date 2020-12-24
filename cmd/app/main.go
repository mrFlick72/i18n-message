package main

import (
	"github/mrflick72/i18n-message/configuration"
	"github/mrflick72/i18n-message/internal/web"
	"sync"
)

func main() {
	// Creates an iris application without any middleware by default
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go configuration.DocumentUpdatesListener(wg)
	go web.NewApplicationServer(wg)
	go web.NewActuatorServer(wg)

	wg.Wait()
}
