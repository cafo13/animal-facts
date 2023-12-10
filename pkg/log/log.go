package log

import (
	"os"
	"time"

	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	log.Logger().Info("logger enabled")
}
