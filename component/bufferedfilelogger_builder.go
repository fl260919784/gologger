package component

import (
	"github.com/fl260919784/gologger/writer"
)

type BufferedFileLogerBuilder struct {
	bufferSize   int
	filename     string
	maxFileSize  uint64
	maxFileCount uint16
}

func (bfb *BufferedFileLogerBuilder) SetBufferSize(bufferSize int) {
	bfb.bufferSize = bufferSize
}

func (bfb *BufferedFileLogerBuilder) SetMaxFileSize(maxFileSize uint64) {
	bfb.maxFileSize = maxFileSize
}

func (bfb *BufferedFileLogerBuilder) SetMaxFileCount(maxFileCount uint16) {
	bfb.maxFileCount = maxFileCount
}

func (bfb *BufferedFileLogerBuilder) SetFilename(filename string) {
	bfb.filename = filename
}

func (bfb *BufferedFileLogerBuilder) Build() Writer {
	if bfb.bufferSize == 0 || len(bfb.filename) == 0 {
		return nil
	}

	var bufferedfileloger writer.WriterDecorator = nil

	bufferBuilder := NewFileBufferWriterDecoratorBuilder()
	bufferBuilder.SetBufferSize(bfb.bufferSize)
	bufferedfileloger = bufferBuilder.Build()
	if bufferedfileloger == nil {
		return nil
	}

	builder := RotateFileWriterBuilder{}
	builder.SetFilename(bfb.filename)
	builder.SetMaxFileSize(bfb.maxFileSize)
	builder.SetMaxFileCount(bfb.maxFileCount)
	writer := builder.Build()
	if writer == nil {
		return nil
	}

	bufferedfileloger.Wrap(writer)

	return bufferedfileloger
}
