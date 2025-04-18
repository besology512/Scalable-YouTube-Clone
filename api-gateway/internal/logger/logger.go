package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})
	log.SetLevel(logrus.InfoLevel)
}
