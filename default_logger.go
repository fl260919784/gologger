package gologger

import (
	"github.com/fl260919784/gologger/writer"
)

var (
	defaultLogger *Logger = NewLogger()
)

func initialize() {
	builder := writer.NewSimpleFileWriterDecoratorBuilder()
	builder.SetFilename("/dev/stdout")
	var w writer.Writer = builder.Build()
	if w == nil {
		w = writer.NewNullWriter()
	}

	defaultLogger.SetLevel(INFO)
	defaultLogger.SetWriter(w)
}

func init() {
	initialize()
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func DefaultLogger() *Logger {
	return defaultLogger
}
