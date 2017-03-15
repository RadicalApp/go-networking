package logger

import (
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
)

var Logger = newLogger()
var Registered = make(map[string]*logging.Logger)
var Backend = newBackend()

//var Available = []string{"git.cyberdust.com/cyberdust-go/dust-client-core/service/connection"}

var Available = []string{}

var format = logging.MustStringFormatter(
	//`%{color}[%{level:.4s}] %{time:15:04:05.000} [%{shortpkg}] %{shortfile}  ▶ %{id:03x}%{color:reset} %{message}`,
	`%{color}[%{level}] %{time:15:04:05.000} [%{shortpkg}] %{shortfile}  ▶ %{color:reset} %{message}`,
)

// Register creates a new logger if this component is set for logging
func Register(name string) *logging.Logger {
	// Create logger
	logger := &logging.Logger{}

	// Add this logger to our map of loggers
	Registered[name] = logger

	return logger
}

// Set up checks a registered logger and adds the appropriate backend.
func setup(name string, logger *logging.Logger) {
	shouldLog := false

	// Check to see if name should be logged.
	for _, l := range Available {
		if l == name {
			shouldLog = true
		}

		// If "all" is in our list of available logging, enable all logging.
		if l == "all" {
			shouldLog = true
		}
	}

	// Set up backend if we should log.
	if shouldLog {
		logger.SetBackend(Backend)
	} else {
		nullBackend := logging.NewLogBackend(ioutil.Discard, "", 0)
		backendLeveled := logging.AddModuleLevel(nullBackend)
		backendLeveled.SetLevel(logging.CRITICAL, "")
		logger.SetBackend(backendLeveled)
	}
}

func newLogger() *logging.Logger {
	logger := logging.MustGetLogger("dust")
	backend := newBackend()
	logger.SetBackend(backend)

	return logger
}

func newBackend() logging.LeveledBackend {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)

	return backendLeveled
}

// Watch will add the given logger to the list of available loggers.
func Watch(logger string) {
	Available = append(Available, logger)
}

// Init initializes logging. This should ALWAYS be called in your program.
func Init(level int) {
	// Convert log level int to logging level type.
	var logLevel logging.Level
	switch level {
	case 1:
		logLevel = logging.DEBUG
	case 2:
		logLevel = logging.INFO
	case 3:
		logLevel = logging.WARNING
	case 4:
		logLevel = logging.ERROR
	}

	// Set up our available loggers.
	for name, logger := range Registered {
		setup(name, logger)
	}

	// Set our logging level.
	Backend.SetLevel(logLevel, "")

	// Set the backends to be used.
	Logger.SetBackend(Backend)

	//Logger.Info("Logging initialized")
}
