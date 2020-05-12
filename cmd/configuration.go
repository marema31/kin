package cmd

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var log = logrus.New()
var logger = log.WithField("app", "kin")

func configureLogging() {
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logger.Errorf("Error getting level: %v", err)
	}

	log.SetLevel(level)

	logFile := viper.GetString("log.path")

	// configure log format
	var formatter logrus.Formatter
	if viper.GetBool("log.json") {
		formatter = &logrus.JSONFormatter{}
	} else {
		disableColors := len(logFile) > 0 && logFile != "-"
		formatter = &logrus.TextFormatter{DisableColors: disableColors, FullTimestamp: true, DisableSorting: true}
	}

	log.SetFormatter(formatter)

	if len(logFile) > 0 && logFile != "-" {
		dir := filepath.Dir(logFile)

		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Errorf("Failed to create log path %s: %s", dir, err)
		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			logger.Errorf("Error while opening log file %s: %v", logFile, err)
		} else {
			log.Out = file
		}

		logrus.RegisterExitHandler(func() {
			if err := file.Close(); err != nil {
				logger.Errorf("Error while closing log: %v", err)
			}
		})
	}
}
