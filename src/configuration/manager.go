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
		configurationManager.init()
	})

	return configurationManager
}
func (manager *Manager) init() {
	manager.viper = viper.New()

	manager.viper.SetConfigName(os.Getenv("CONFIGURATION_FILE_NAME"))
	manager.viper.SetConfigType(os.Getenv("CONFIGURATION_FILE_TYPE"))
	manager.viper.AddConfigPath(os.Getenv("CONFIGURATION_PATH"))
	manager.viper.ReadInConfig()
}

func (manager *Manager) GetConfigFor(configKey string) string {
	return manager.viper.GetString(configKey)
}

func (manager *Manager) GetStringMapFor(configKey string) map[string]string {
	return manager.viper.GetStringMapString(configKey)
}
