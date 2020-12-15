package logging

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

// InitUserLogger sets up the logger which writes to the DB (dependent on level as specified in the logging hook)
func InitUserLogger(logLevel, logfile string) {
	logFile, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	// Show where error was logged, function, line number, etc.
	logrus.SetReportCaller(true)

	// Output to stdout and logfile
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Set the formatter for the logging to make them more human readable
	logrus.SetFormatter(&nested.Formatter{})

	// Only log the warning severity or above.
	logrus.SetLevel(parseLogLevelFromConfig(logLevel))
}

// InitInternalLogger sets up the system logger which only writes to the console,
// most errors are written to this logger
func InitInternalLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Caller().
		Logger()
}

func parseLogLevelFromConfig(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
