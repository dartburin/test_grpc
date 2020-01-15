package logger

import (
	"os"

	lg "github.com/sirupsen/logrus"
)

func LogInit(lvl string, fName string) lg.FieldLogger {
	log := lg.New()

	switch lvl {
	case "Trace":
		log.Level = lg.TraceLevel
	case "Debug":
		log.Level = lg.DebugLevel
	case "Info":
		log.Level = lg.InfoLevel
	case "Warn":
		log.Level = lg.WarnLevel
	default:
		log.Level = lg.DebugLevel //lg.ErrorLevel
	}

	fmtTime := new(lg.TextFormatter)
	fmtTime.TimestampFormat = "2006-01-02 15:04:05"
	fmtTime.FullTimestamp = true
	log.SetFormatter(fmtTime)

	if fName != "" {
		//move outut to file
		fl, err := os.OpenFile(fName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Errorf("Don`t open %s file.", fName)
		} else {
			log.SetOutput(fl)
		}
	}
	return log
}

func PrintOsArgs(log lg.FieldLogger) {
	for _, a := range os.Args {
		log.Infof("%s ", a)
	}
	log.Infoln("")
}
