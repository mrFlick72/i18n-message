package main

import (
	"github/mrflick72/i18n-message/configuration"
	"sync"
)

func main() {
	// Creates an iris application without any middleware by default
	configurationWg := &sync.WaitGroup{}
	configurationWg.Add(1)
	manager := configuration.GetConfigurationManagerInstance()

	go manager.Init(configurationWg)
	configurationWg.Wait()

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go configuration.DocumentUpdatesListener(wg)
	go configuration.NewApplicationServer(wg)
	go configuration.NewActuatorServer(wg)

	wg.Wait()
}
