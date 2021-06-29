package configuration

import (
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	configurationManager *Manager
	once                 sync.Once
)

type Manager struct {
	viper *viper.Viper
}

func GetConfigurationManagerInstance() *Manager {
	once.Do(func() {
		configurationManager = &Manager{}
	})

	return configurationManager
}
func (manager *Manager) Init(wg *sync.WaitGroup) {
	manager.viper = viper.New()

	manager.viper.SetConfigName(os.Getenv("CONFIGURATION_FILE_NAME"))
	manager.viper.SetConfigType(os.Getenv("CONFIGURATION_FILE_TYPE"))
	manager.viper.AddConfigPath(os.Getenv("CONFIGURATION_PATH"))
	manager.viper.ReadInConfig()
	wg.Done()
}

func (manager *Manager) GetConfigFor(configKey string) string {
	return manager.viper.GetString(configKey)
}
