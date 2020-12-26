package main

import (
	"github/mrflick72/i18n-message/src/configuration"
	"github/mrflick72/i18n-message/src/configuration/application"
	"sync"
)

func main() {
	initConfigurationManager()
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

func initConfigurationManager() {
	configurationWg := &sync.WaitGroup{}
	configurationWg.Add(1)
	manager := configuration.GetConfigurationManagerInstance()

	go manager.Init(configurationWg)
	configurationWg.Wait()
}
