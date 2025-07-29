package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Level
type (
	Level int

	color string

	entry struct {
		Name       string `json:"name,omitempty"`
		Level      string `json:"level"`
		Caller     string `json:"caller,omitempty"`
		Message    string `json:"msg,omitempty"`
		Fields     []any  `json:"fields,omitempty"`
		Stacktrace string `json:"stacktrace,omitempty"`
		Timestamp  string `json:"timestamp"`
	}
)

const (
	// DEBUG
	DEBUG Level = iota
	// INFO
	INFO
	// WARN
	WARN
	// ERROR
	ERROR
	// FATAL
	FATAL

	colorReset  color = "\033[0m"
	colorRed    color = "\033[31m"
	colorGreen  color = "\033[32m"
	colorYellow color = "\033[33m"
	colorBlue   color = "\033[34m"
	colorPurple color = "\033[35m"
	colorCyan   color = "\033[36m"
	colorWhite  color = "\033[37m"
)

func (lvl Level) string() string {
	switch lvl {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

func levelFromString(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return DEBUG
	}
}

func colorFromLevel(level Level) color {
	switch level {
	case DEBUG:
		return colorCyan
	case INFO:
		return colorGreen
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	case FATAL:
		return colorPurple
	default:
		return colorReset
	}
}

// LoggerOption
type LoggerOption func(*Logger)

// WithWriter
func WithWriter(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.w = w
	}
}

// WithName
func WithName(n string) LoggerOption {
	return func(l *Logger) {
		l.name = n
	}
}

// WithLevel
func WithLevel(level string) LoggerOption {
	return func(l *Logger) {
		l.minLevel = levelFromString(level)
	}
}

// WithCaller
func WithCaller(includeCaller bool) LoggerOption {
	return func(l *Logger) {
		l.includeCaller = includeCaller
	}
}

// WithStacktrace
func WithStacktrace(stacktrace bool) LoggerOption {
	return func(l *Logger) {
		l.stacktrace = stacktrace
	}
}

// Logger
type Logger struct {
	mu sync.Mutex
	w  io.Writer

	minLevel Level

	name string

	includeCaller bool
	stacktrace    bool
}

// New
func New(opts ...LoggerOption) *Logger {
	l := &Logger{
		mu: sync.Mutex{},
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Debug
func (l *Logger) Debug(format string) {
	l.logMsg(DEBUG, "%s", format)
}

// Debugf
func (l *Logger) Debugf(format string, args ...any) {
	l.logMsg(DEBUG, format, args...)
}

// Info
func (l *Logger) Info(format string) {
	l.logMsg(INFO, "%s", format)
}

// Infof
func (l *Logger) Infof(format string, args ...any) {
	l.logMsg(INFO, format, args...)
}

// Warn
func (l *Logger) Warn(format string) {
	l.logMsg(WARN, "%s", format)
}

// Warnf
func (l *Logger) Warnf(format string, args ...any) {
	l.logMsg(WARN, format, args...)
}

// Error
func (l *Logger) Error(format string) {
	l.logMsg(ERROR, "%s", format)
}

// Errorf
func (l *Logger) Errorf(format string, args ...any) {
	l.logMsg(ERROR, format, args...)
}

// Fatal
func (l *Logger) Fatal(format string) {
	l.logMsg(FATAL, "%s", format)
}

// Fatalf
func (l *Logger) Fatalf(format string, args ...any) {
	l.logMsg(FATAL, format, args...)
}

func (l *Logger) logMsg(level Level, format string, args ...any) {
	// level is greater than min
	if l.minLevel > level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	logEntry := l.fmt(level, format, args...)
	entrybytes, err := json.Marshal(&logEntry)
	if err != nil {
		fmt.Fprintf(l.w, "failed to marshal message %v original message %s", err, format) // nolint: errcheck
	}

	entrybytes = append(entrybytes, '\n')
	_, err = l.w.Write(entrybytes)
	if err != nil {
		fmt.Fprintf(l.w, "failed to write to out: %v", err) // nolint: errcheck
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) fmt(level Level, format string, args ...any) entry {
	ts := time.Now().Format("2006-01-02 15:04:05")

	var en entry
	en.Level = level.string()
	en.Message = format
	en.Fields = make([]any, 0, len(args))
	en.Name = l.name
	en.Timestamp = ts

	for _, arg := range args {
		en.Fields = append(en.Fields, arg)
		if level == ERROR || level == FATAL {
			if l.stacktrace {
				stacktrace := debug.Stack()
				en.Stacktrace = fmt.Sprintf(" %sstracktrace %s%s", colorFromLevel(level), colorReset, string(stacktrace))
			}
		}
	}

	return en
}
