package gologger

import (
	"fmt"
	"github.com/fl260919784/gologger/writer"
	"time"
)

const (
	ERROR = iota
	WARN
	INFO
	DEBUG
)

var levels = [...]string{
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
}

func String2Level(level string) int {
	for index, name := range levels {
		if name == level {
			return index
		}
	}

	return 0
}

type Logger struct {
	level  int
	writer writer.Writer
}

func NewLogger() *Logger {
	return &Logger{
		level:  INFO,
		writer: writer.NewNullWriter(),
	}
}

func (l *Logger) prefix(level int) string {
	now := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%03d [%s] ",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(),
		now.Second(), now.Nanosecond()/1000000, levels[level])
}

func (l *Logger) suffix() string {
	return ""
}

// 保持每条日志最后只有一个换行符
func (l *Logger) formatter(msg string) string {
	var lastIndex int = len(msg) - 1
	for lastIndex >= 0 {
		if msg[lastIndex] != '\n' {
			break
		}

		lastIndex--
	}

	return fmt.Sprint(msg[0:lastIndex+1], "\n")
}

func (l *Logger) SetWriter(writer writer.Writer) {
	l.writer = writer
}

func (l *Logger) SetLevel(level int) bool {
	if level < ERROR || level > DEBUG {
		return false
	}

	l.level = level
	return true
}

func (l *Logger) Flush() {
	l.writer.Flush()
}

func (l *Logger) Close() {
	l.writer.Close()
}

func (l *Logger) output(level int, format string, v ...interface{}) {
	if level > l.level {
		return
	}

	l.writer.Write(l.prefix(level), l.formatter(fmt.Sprintf(format, v...)), l.suffix())
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(DEBUG, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(INFO, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(WARN, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(ERROR, format, v...)
}
