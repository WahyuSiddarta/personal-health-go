package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/config"
	database "github.com/WahyuSiddarta/be_saham_go/db"
	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/router"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"golang.org/x/term"
	"gopkg.in/natefinch/lumberjack.v2"
)

const tstamp = "2006-01-02 15:04:05"

// LogConfig : Configuration for logging
type LogConfig struct {
	// Enable console logging
	ConsoleLoggingEnabled bool
	// EncodeLogsAsJSON makes the log framework log JSON
	EncodeLogsAsJSON bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

// newRollingFile creates a new rolling file logger
func NewRollingFile(config LogConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		fmt.Printf("can't create log directory: %v", err)
		return nil
	}

	logFile := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}

	return logFile
}

// initLogger initializes the logging system
func InitLogger() *zerolog.Logger {
	config := LogConfig{
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJSON:      true,
		FileLoggingEnabled:    true,
		Directory:             "./log",
		Filename:              "backend.log",
		MaxSize:               100,
		MaxBackups:            7,
		MaxAge:                30,
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	var writers []io.Writer

	// Only add console writer if console logging is enabled
	if config.ConsoleLoggingEnabled {
		if term.IsTerminal(int(os.Stdout.Fd())) {
			if runtime.GOOS == "windows" {
				writers = append(writers, zerolog.ConsoleWriter{
					Out:        colorable.NewColorableStdout(),
					TimeFormat: tstamp,
				})
			} else {
				writers = append(writers, zerolog.ConsoleWriter{
					Out:        os.Stdout,
					TimeFormat: tstamp,
					FormatTimestamp: func(i interface{}) string {
						parse, _ := time.Parse(time.RFC3339Nano, i.(string))
						x, _ := helper.TimeInWIB(parse)
						return "\033[1;36m" + x.Format(tstamp) + "\033[0m"
					},
				})
			}
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: tstamp,
			})
		}
	}

	// Add file writer if file logging is enabled
	if config.FileLoggingEnabled {
		fileWriter := NewRollingFile(config)
		if fileWriter != nil {
			writers = append(writers, fileWriter)
		}
	}

	// Create multi-writer only if we have writers
	if len(writers) == 0 {
		// Fallback to stdout if no writers configured
		writers = append(writers, os.Stdout)
	}

	mw := io.MultiWriter(writers...)

	// Setup diode writer for async logging
	wr := diode.NewWriter(mw, 500, 50*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})

	// Create logger with or without JSON encoding based on config
	var logger zerolog.Logger
	if config.EncodeLogsAsJSON {
		logger = zerolog.New(wr).With().Timestamp().Logger()
	} else {
		// Only use non-JSON encoding for console output
		// For file output, always use JSON
		if config.FileLoggingEnabled && config.ConsoleLoggingEnabled {
			// Mixed logging - console gets formatted, file gets JSON
			consoleWriter := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: tstamp,
			}

			fileWriter := NewRollingFile(config)
			combinedWriter := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
			logger = zerolog.New(combinedWriter).With().Timestamp().Logger()
		} else {
			// Single output mode
			logger = zerolog.New(wr).With().Timestamp().Logger()
		}
	}

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJSON).
		Bool("consoleLogging", config.ConsoleLoggingEnabled).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("Logging system configured")

	return &logger
}

func DistrubuteLogger(logger *zerolog.Logger) {
	helper.Logger = logger
	database.Logger = logger
	models.Logger = logger
	api.Logger = logger
	config.Logger = logger
	router.Logger = logger
	middleware.Logger = logger
}
