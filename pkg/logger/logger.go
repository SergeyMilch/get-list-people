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

func Info(msg string) {
	log.Info(msg)
}

func Warn(msg string) {
	log.Warn(msg)
}

func Error(msg string) {
	log.Error(msg)
}
