package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})
	Log.SetLevel(logrus.InfoLevel)
}
