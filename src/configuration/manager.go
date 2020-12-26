package configuration

import (
	"github.com/Piszmog/cloudconfigclient"
	"github.com/spf13/viper"
	"github/mrflick72/i18n-message/src/internal/logging"
	"net/http"
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
	configClient, err := cloudconfigclient.NewLocalClient(&http.Client{}, []string{os.Getenv("spring.cloud.config.uri")})
	config, err := configClient.GetConfiguration("i18n-messages", []string{os.Getenv("spring.profiles.active")})
	if err != nil {
		panic(err)
	}

	manager.viper = viper.New()
	for key, value := range config.PropertySources[1].Source {
		manager.viper.SetDefault(key, value)
	}

	viper.SetConfigName("application.yaml")
	viper.SetConfigType("yaml")
	manager.viper.AddConfigPath(".")
	manager.viper.AutomaticEnv()
	manager.viper.ReadInConfig()

	logging.LogDebugFor(manager.viper)
	wg.Done()
}

func (manager *Manager) GetConfigFor(configKey string) string {
	return manager.viper.GetString(configKey)
}
func (manager *Manager) GetStringMapFor(configKey string) map[string]string {
	return manager.viper.GetStringMapString(configKey)
}
