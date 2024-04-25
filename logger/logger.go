package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	l    *zerolog.Logger
	once sync.Once
)

// Initialize configures and initializes the global logger.
// It sets up a console writer and a rolling file writer with specified rotation policies.
func Initialize() *zerolog.Logger {
	once.Do(func() {
		// consoleWriter formats logs with a human-friendly, colorized output.
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		// fileWriter writes logs to a file with log rotation based on size,
		// number of backups, and age of the log files.
		fileWriter := &lumberjack.Logger{
			Filename:   "log/app.log",
			MaxSize:    5,    // MaxSize is the maximum size in MB before the log is rotated.
			MaxBackups: 10,   // MaxBackups is the maximum number of log files to keep.
			MaxAge:     14,   // MaxAge is the maximum number of days to keep a log file.
			Compress:   true, // Compress indicates if the rotated log files should be compressed (gzip).
		}

		// MultiLevelWriter allows logs to be written to both console and file writers.
		output := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

		// The logger is configured with the default DebugLevel log level
		// and a timestamp is added to every log message.
		logger := zerolog.New(output).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		l = &logger
	})
	return l
}

// Get returns the global logger instance. It ensures the logger is initialized
// only once during the application lifetime.
func Get() *zerolog.Logger {
	return Initialize()
}

// logWithContext logs a message with the context information (like request ID) using the global logger.
func logWithContext(ctx echo.Context, level zerolog.Level, msg string, err error) {
	// Extract the request ID from the echo context.
	requestID := ctx.Request().Header.Get(echo.HeaderXRequestID)

	// The logger instance is enriched with the request ID.
	logger := Get().With().Str("request_id", requestID).Logger()

	// Logs the message with an error if provided, otherwise with the specified level.
	if err != nil {
		logger.Error().Err(err).Msg(msg)
	} else {
		logger.WithLevel(level).Msg(msg)
	}
}

// The following functions wrap logWithContext providing a convenient way
// to log messages with different log levels: Info, Debug, Warn, and Error.

func Infof(ctx echo.Context, msg string, v ...interface{}) {
	logWithContext(ctx, zerolog.InfoLevel, fmt.Sprintf(msg, v...), nil)
}

func Debugf(ctx echo.Context, msg string, v ...interface{}) {
	logWithContext(ctx, zerolog.DebugLevel, fmt.Sprintf(msg, v...), nil)
}

func Warnf(ctx echo.Context, msg string, v ...interface{}) {
	logWithContext(ctx, zerolog.WarnLevel, fmt.Sprintf(msg, v...), nil)
}

func Errorf(ctx echo.Context, err error, msg string, v ...interface{}) {
	logWithContext(ctx, zerolog.ErrorLevel, fmt.Sprintf(msg, v...), err)
}
