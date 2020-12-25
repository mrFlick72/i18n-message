package configuration

import (
	"github.com/Piszmog/cloudconfigclient"
	"github.com/spf13/viper"
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
	configClient, err := cloudconfigclient.NewLocalClient(&http.Client{}, []string{os.Getenv("CONFIG_SERVER_URL")})
	config, err := configClient.GetConfiguration("i18n-messages", []string{os.Getenv("APPLICATION_PROFILE")})
	if err != nil {
		panic(err)
	}

	manager.viper = viper.New()
	for key, value := range config.PropertySources[0].Source {
		manager.viper.SetDefault(key, value)
	}

	viper.SetConfigName("application.yaml")
	viper.SetConfigType("yaml")
	manager.viper.AddConfigPath(".")
	manager.viper.AutomaticEnv()
	manager.viper.ReadInConfig()

	wg.Done()
}

func (manager *Manager) GetConfigFor(configKey string) string {
	return manager.viper.GetString(configKey)
}
func (manager *Manager) GetStringMapFor(configKey string) map[string]string {
	return manager.viper.GetStringMapString(configKey)
}
