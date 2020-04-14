package logger

import (
	"github.com/fl260919784/gologger/writer"
)

var (
	defaultAccessLogger *Logger = NewLogger()
)

func initializeAccessLogger() {
	builder := writer.NewSimpleFileWriterDecoratorBuilder()
	builder.SetFilename("/dev/stdout")
	w := factory.Build()
	if w == nil {
		w = NewNullWriter()
	}

	defaultAccessLogger.SetLevel(INFO)
	defaultAccessLogger.SetWriter(w)
}

func init() {
	initializeAccessLogger()
}

func AccessLogger() *Logger {
	return defaultAccessLogger
}
