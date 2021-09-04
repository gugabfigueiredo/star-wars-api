package log

import "io"

// convenience code
// mostly taken from https://gist.github.com/panta/2530672ca641d953ae452ecb5ef79d7d

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

// Config - Configuration for logging
type Config struct {
 	Context string `default:"sw-api"`
	// Enable console logging
	ConsoleLoggingEnabled bool `default:"true"`
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool `default:"true"`
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `default:"true"`
	// Directory to log to to when filelogging is enabled
	Directory string `default:"/var/log"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `default:"sw-api-server"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `default:"10"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `default:"10"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `default:"2"`
}

// New sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func New(config *Config) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Str("context", config.Context).
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	return &Logger{
		Logger: &logger,
		context: map[string]interface{}{"context": config.Context},
	}
}

func newRollingFile(config *Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Context, config.Filename),
		MaxBackups: config.MaxBackups,		// files
		MaxSize:    config.MaxSize,			// megabytes
		MaxAge:     config.MaxAge,			// days
	}
}