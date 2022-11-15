package main

import (
	"github/mrflick72/i18n-message/src/configuration/application"
	"sync"
)

func main() {
	initApplicationServer()
}

func initApplicationServer() {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go application.DocumentUpdatesListener(wg)
	go application.NewApplicationServer(wg)
	go application.NewActuatorServer(wg)
	wg.Wait()
}
