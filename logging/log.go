// Package logging provides structured logging with logrus.
package logging

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	// Logger is a configured logrus.Logger.
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
}

func CreateLog() {

	file, err := os.OpenFile("/Users/applemd1011/Desktop/Go-Projects/logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		Log.Fatal(err)
	}
	//defer file.Close()
	Log.SetOutput(file)
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}
