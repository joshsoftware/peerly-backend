package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	l "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *l.Logger

var lumberjackLogger = &lumberjack.Logger{
	MaxSize:    100, // megabytes
	MaxBackups: 3,   // number of log files
	MaxAge:     365, // days
	Compress:   true,
}

// Error - prints out an error
func Error(ctx context.Context, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Error(args...)
}

// Errorf - prints out an error with formatted output
func Errorf(ctx context.Context, format string, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Errorf(format, args...)
}

// Warn - prints out a warning
func Warn(ctx context.Context, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Warn(args...)
}

// Fatal - will print out the error info and exit the program
func Fatal(ctx context.Context, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Fatal(args...)
}

// Info - prints out basic information
func Info(ctx context.Context, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Info(args...)
}

// Infof - prints out basic information
func Infof(ctx context.Context, format string, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Infof(format, args...)
}

// Debug - prints out debug information
func Debug(ctx context.Context, args ...interface{}) {
	log := getLoggerWithRequestContext(ctx)
	log.Debug(args...)
}

func SetupLogger() (*l.Logger, error) {

	lumberjackLogger.Filename = fmt.Sprintf("/var/log/peerly/%s_peerly_backend.log", time.Now().Format("2006-01-02_15-04-05"))
	file, err := os.Create(lumberjackLogger.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}
	file.Close()

	// Initialize Logrus logger
	logger := l.New()
	logger.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))
	logger.SetFormatter(&l.TextFormatter{
		FullTimestamp: true,
	})

	// Set the logging level
	logger.SetLevel(l.InfoLevel)

	Logger = logger
	return logger, nil
}

func getLoggerWithRequestContext(ctx context.Context) *l.Entry {
	if Logger == nil {
		Logger = l.New()
	}

	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = "N/A"
	}

	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		userID = 0
	}

	return Logger.WithFields(l.Fields{
		"req_id":  requestID,
		"user_id": userID,
	})
}

// Cron Logger Implementation
var CronLogger *l.Logger

var cronLumberjackLogger = &lumberjack.Logger{
	MaxSize:    100, // megabytes
	MaxBackups: 3,   // number of log files
	MaxAge:     365, // days
	Compress:   true,
}

func SetupCronLogger() (*l.Logger, error) {
	logDir := "/var/log/peerly"
	logFilename := "cron.log"
	cronLumberjackLogger.Filename = filepath.Join(logDir, logFilename)

	if err := os.MkdirAll(filepath.Dir(cronLumberjackLogger.Filename), 0755); err != nil {
		logDir = "./logs"
		cronLumberjackLogger.Filename = filepath.Join(logDir, logFilename)
		if err := os.MkdirAll(filepath.Dir(cronLumberjackLogger.Filename), 0755); err != nil {
			return nil, fmt.Errorf("failed to create cron log directory: %w", err)
		}
	}

	file, err := os.OpenFile(cronLumberjackLogger.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create cron log file: %w", err)
	}
	file.Close()

	// Initialize Logrus logger
	logger := l.New()
	logger.SetOutput(io.MultiWriter(os.Stdout, cronLumberjackLogger))
	logger.SetFormatter(&l.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetLevel(l.InfoLevel)

	CronLogger = logger
	return logger, nil
}

func getCronLoggerWithRequestContext(ctx context.Context) *l.Entry {
	if CronLogger == nil {
		CronLogger = l.New()
	}

	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = "N/A"
	}

	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		userID = 0
	}

	return CronLogger.WithFields(l.Fields{
		"req_id":  requestID,
		"user_id": userID,
	})
}

func CronError(ctx context.Context, args ...interface{}) {
	log := getCronLoggerWithRequestContext(ctx)
	log.Error(args...)
}

func CronErrorf(ctx context.Context, format string, args ...interface{}) {
	log := getCronLoggerWithRequestContext(ctx)
	log.Errorf(format, args...)
}

func CronWarn(ctx context.Context, args ...interface{}) {
	log := getCronLoggerWithRequestContext(ctx)
	log.Warn(args...)
}

func CronInfo(ctx context.Context, args ...interface{}) {
	log := getCronLoggerWithRequestContext(ctx)
	log.Info(args...)
}

func CronInfof(ctx context.Context, format string, args ...interface{}) {
	log := getCronLoggerWithRequestContext(ctx)
	log.Infof(format, args...)
}

func CronDebug(ctx context.Context, args ...interface{}) {
	log := getCronLoggerWithRequestContext(ctx)
	log.Debug(args...)
}

