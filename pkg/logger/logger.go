package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	log.Level = logrus.DebugLevel
}

func Error(errMsg string, details string) {
	log.WithFields(logrus.Fields{
		"details": details,
	}).Error(errMsg)
}

func Info(errMsg string, details string) {
	log.WithFields(logrus.Fields{
		"details": details,
	}).Error(errMsg)
}

func Warn(errMsg string, details string) {
	log.WithFields(logrus.Fields{
		"details": details,
	}).Error(errMsg)
}
