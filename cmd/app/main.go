package main

import (
	"github/mrflick72/i18n-message/configuration"
	"sync"
)

func main() {
	// Creates an iris application without any middleware by default
	initConfigurationManager()
	initApplicationServer()
}

func initApplicationServer() {
	wg := &sync.WaitGroup{}
	wg.Add(4)
	go configuration.DocumentUpdatesListener(wg)
	go configuration.NewApplicationServer(wg)
	go configuration.NewActuatorServer(wg)
	go configuration.NewARSocketServer(wg)
	wg.Wait()
}

func initConfigurationManager() {
	configurationWg := &sync.WaitGroup{}
	configurationWg.Add(1)
	manager := configuration.GetConfigurationManagerInstance()

	go manager.Init(configurationWg)
	configurationWg.Wait()
}
