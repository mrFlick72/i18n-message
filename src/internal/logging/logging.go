package logging

import (
	"fmt"
	"github/mrflick72/i18n-message/src/configuration"
	"go.uber.org/zap"
	"os"
)

var (
	manager   = configuration.GetConfigurationManagerInstance()
	fileName  = manager.GetConfigFor("logging.fileName")
	logger, _ = loggerConfigurer().Build()
)

func loggerConfigurer() zap.Config {
	cfg := zap.NewProductionConfig()

	fmt.Println("log file name: ", fileName)
	if len(fileName) > 0 {
		_, err := os.Create(fileName)
		if err != nil {
			panic("log file des not exist")
		}

		cfg.OutputPaths = []string{fileName}
	}
	return cfg
}

func Logger() *zap.Logger {
	return logger
}

func LogErrorFor(message interface{}) {
	str := fmt.Sprintf("%v", message)
	logger.Error(str)
}

func LogInfoFor(message interface{}) {
	str := fmt.Sprintf("%v", message)
	logger.Info(str)
}
func LogDebugFor(message interface{}) {
	str := fmt.Sprintf("%v", message)
	logger.Debug(str)
}

func Dispose() {
	panic("TODO it have to be implemented")
}
