package log

import (
	"context"
	"fmt"
	"io"
	"os"
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
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Error(requestID, args)
}

// Errorf - prints out an error with formatted output
func Errorf(ctx context.Context, format string, args ...interface{}) {
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Errorf("req_id: [%s] "+format, append([]interface{}{requestID}, args...)...)
}

// Warn - prints out a warning
func Warn(ctx context.Context, args ...interface{}) {
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Warn("req_id: ", requestID, args)
}

// Fatal - will print out the error info and exit the program
func Fatal(ctx context.Context, args ...interface{}) {
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Fatal("req_id: ", requestID, args)
}

// Info - prints out basic information
func Info(ctx context.Context, args ...interface{}) {
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Info("req_id: ", requestID, args)
}

// Infof - prints out basic information
func Infof(ctx context.Context, format string, args ...interface{}) {
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Infof("req_id: [%s] "+format, append([]interface{}{requestID}, args...)...)
}

// Debug - prints out debug information
func Debug(ctx context.Context, args ...interface{}) {
	requestID, ok := ctx.Value(constants.RequestID).(string)
	if !ok {
		requestID = ""
	}
	Logger.Debug("req_id: ", requestID, args)
}

func SetupLogger() (*l.Logger, error) {

	lumberjackLogger.Filename = fmt.Sprintf("./logs/%s_peerly_backend.log", time.Now().Format("2006-01-02_15-04-05"))
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
