package log

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Level
type Level int

type color string

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

// Logger
type Logger struct {
	mu sync.Mutex
	w  io.Writer

	minLevel Level

	name string
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
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logMsg(DEBUG, format, args...)
}

// Info
func (l *Logger) Info(format string) {
	l.logMsg(INFO, "%s", format)
}

// Infof
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logMsg(INFO, format, args...)
}

// Warn
func (l *Logger) Warn(format string) {
	l.logMsg(WARN, "%s", format)
}

// Warnf
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logMsg(WARN, format, args...)
}

// Error
func (l *Logger) Error(format string) {
	l.logMsg(ERROR, "%s", format)
}

// Errorf
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logMsg(ERROR, format, args...)
}

// Fatal
func (l *Logger) Fatal(format string) {
	l.logMsg(FATAL, "%s", format)
}

// Fatalf
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logMsg(FATAL, format, args...)
}

func (l *Logger) logMsg(level Level, format string, args ...interface{}) {
	// level is greater than min
	if l.minLevel > level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	fmsg := l.fmt(level, format, args...)
	_, err := l.w.Write([]byte(fmsg))
	if err != nil {
		fmt.Fprintf(l.w, "failed to write to out: %v", err)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) fmt(level Level, format string, args ...interface{}) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	var clr color

	if l.w == os.Stdout || l.w == os.Stdin {
		clr = colorFromLevel(level)
	}

	lmsg := fmt.Sprintf("%s%s [%s%s%s] %s %s", clr, ts, colorFromLevel(level), level.string(), colorReset, l.name, msg)

	if level == ERROR || level == FATAL {
		stacktrace := debug.Stack()
		lmsg += fmt.Sprintf(" %sstracktrace %s%s", colorFromLevel(level), colorReset, string(stacktrace))
	}

	lmsg += "\n"

	return lmsg
}
