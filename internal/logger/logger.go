package logger

import (
  "fmt"
  "os"
  "log/slog"
)


var defaultLogger Logger

type Logger interface {
  // Info logs an informational message
  LogInfo(msg string, args ...any)
  // Error logs an error message
  LogError(msg string, args ...any)
  // Debug logs a debug message
  LogDebug(msg string, args ...any)

  // Close closes the logger
  Close() error
}

type LoggerClient struct {
  logger *slog.Logger
  file   *os.File
}

func (l *LoggerClient) LogInfo(msg string, args ...any) {
  l.logger.Info(msg, args...)
}

func (l *LoggerClient) LogError(msg string, args ...any) {
  l.logger.Error(msg, args...)
}

func (l *LoggerClient) LogDebug(msg string, args ...any) {
  l.logger.Debug(msg, args...)
}

func (l *LoggerClient) Close() error {
  if l.file != nil {
    err := l.file.Close()
    if err != nil {
      return fmt.Errorf("failed to close log file: %v", err)
    }
  }
  return nil
}

type EmptyLogger struct{}
func (l *EmptyLogger) LogInfo(msg string, args ...any) {
  // No-op
}
func (l *EmptyLogger) LogError(msg string, args ...any) {
  // No-op
}
func (l *EmptyLogger) LogDebug(msg string, args ...any) {
  // No-op
}
func (l *EmptyLogger) Close() error {
  // No-op
  return nil
}

// NewLogger creates a new logger instance
// It accepts a path to a file where logs will be written
// and a boolean indicating whether to log to stdout as well
func NewLogger() (Logger, error) {
  logEnabled := os.Getenv("AEROSPACE_MARKS_LOG")
  if logEnabled == "" {
    return &EmptyLogger{}, nil
  }

  path := os.Getenv("AEROSPACE_MARKS_LOG_FILEPATH")
  if path == "" {
    path = "/tmp/aerospace-marks.log"
  }

  file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
      return nil, fmt.Errorf("failed to open log file: %v", err)
  }

  configLogLevel := os.Getenv("AEROSPACE_MARKS_LOG_LEVEL")
  logLevel := slog.LevelError
  if configLogLevel != "" {
    switch configLogLevel {
    case "DEBUG":
      logLevel = slog.LevelDebug
    case "INFO":
      logLevel = slog.LevelInfo
    case "WARN":
      logLevel = slog.LevelWarn
    default:
      logLevel = slog.LevelError
    }
  }

  textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
      Level: logLevel,
  })

  newLogger := slog.New(textHandler)
  logClient := &LoggerClient{
    logger: newLogger,
    file:   file,
  }

  return logClient, nil
}

func SetDefaultLogger(logger Logger) {
  // Set the default logger to the provided logger
  defaultLogger = logger
}

func GetDefaultLogger() (Logger, error) {
  if defaultLogger == nil {
    return nil, fmt.Errorf("default logger is not set")
  }
  return defaultLogger, nil
}
