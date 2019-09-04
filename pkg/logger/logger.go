package logger

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

type Logger = *log.Entry

// Setup logrus logger.
func Setup(logFormat, logLevel, logFile string) (*log.Logger, error) {
	logger := log.New()

	logger.AddHook(filename.NewHook())

	setFormat(logFormat, logger)
	setOutput(logFile, logger)
	setLevel(logLevel, logger)

	return logger, nil
}

type utcFormatter struct {
	log.Formatter
}

// Format will format the log message with UTC time.
func (u utcFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

func setFormat(logFormat string, logger *log.Logger) {
	switch strings.ToLower(logFormat) {
	case "text":
		logger.SetFormatter(&utcFormatter{&log.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000", FullTimestamp: true}})
	case "json":
		logger.SetFormatter(&utcFormatter{&log.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000"}})
	default:
		logger.Warnf("unknown logformat '%s', falling back to 'text' format.", logFormat)
		logger.SetFormatter(&utcFormatter{&log.TextFormatter{}})
	}
}

func setOutput(logFile string, logger *log.Logger) {
	switch strings.ToLower(logFile) {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "null":
		logger.SetOutput(ioutil.Discard)
	default:
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			_, err = os.Create(logFile)
			if err != nil {
				logger.Warnf("unable to create logfile '%s', falling back to stdout: %v", logFile, err)
				logger.SetOutput(os.Stdout)
				return
			}
		}

		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0640)
		if err != nil {
			logger.Warnf("unable to open logfile '%s', falling back to stdout: %v", logFile, err)
			logger.SetOutput(os.Stdout)
			return
		}
		multiWriter := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(multiWriter)
	}
}

func setLevel(logLevel string, logger *log.Logger) {
	switch strings.ToLower(logLevel) {
	case "debug":
		logger.Info("setting debug level")
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.Info("setting info level")
		logger.SetLevel(log.InfoLevel)
	case "warning":
		logger.Info("setting warning level")
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.Info("setting error level")
		logger.SetLevel(log.ErrorLevel)
	default:
		logger.Warnf("unknown loglevel '%s', falling back to 'info' level.", logLevel)
		logger.SetLevel(log.InfoLevel)
	}
}
