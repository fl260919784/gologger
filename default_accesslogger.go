package gologger

import (
	"github.com/fl260919784/gologger/writer"
)

var (
	defaultAccessLogger *Logger = NewLogger()
)

func initializeAccessLogger() {
	builder := writer.NewSimpleFileWriterDecoratorBuilder()
	builder.SetFilename("/dev/stdout")
	var w writer.Writer = builder.Build()
	if w == nil {
		w = writer.NewNullWriter()
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
