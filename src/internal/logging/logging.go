package logging

import (
	"fmt"
	"github/mrflick72/i18n-message/src/configuration"
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
)

var (
	manager       = configuration.GetConfigurationManagerInstance()
	loggerManager *Logger
	once          sync.Once
)

type Logger struct {
	logger *zap.Logger
}

func GetLoggerInstance() *Logger {
	once.Do(func() {
		loggerManager = new()
	})

	return loggerManager
}

func new() *Logger {
	fileName := manager.GetConfigFor("logging.file.name")
	log.Printf("log fileName: %v", fileName)
	cfg := zap.NewProductionConfig()

	fmt.Println("log file name: ", fileName)
	if len(fileName) > 0 {
		_, err := os.Create(fileName)
		if err != nil {
			panic("log file does not exist")
		}

		cfg.OutputPaths = []string{fileName}
	}
	log, _ := cfg.Build()
	return &Logger{
		logger: log,
	}
}

func (impl *Logger) LogErrorFor(message interface{}) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Error(str)
}

func (impl *Logger) LogInfoFor(message interface{}) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Info(str)
}
func (impl *Logger) LogDebugFor(message interface{}) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Debug(str)
}

func Dispose() {
	panic("TODO it have to be implemented")
}
