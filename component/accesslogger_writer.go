package component

import (
	"github.com/fl260919784/gologger/writer"
)

type AccessloggerBuilder struct {
	raddr string
	rport int
	laddr string
	lport int

	filename     string
	maxFileSize  uint64
	maxFileCount uint16

	bufferSize int
}

func (alb *AccessloggerBuilder) SetRemoteAddr(addr string) {
	alb.raddr = addr
}

func (alb *AccessloggerBuilder) SetRemotePort(port int) {
	alb.rport = port
}

func (alb *AccessloggerBuilder) SetLocalAddr(addr string) {
	alb.laddr = addr
}

func (alb *AccessloggerBuilder) SetLocalPort(port int) {
	alb.lport = port
}

func (alb *AccessloggerBuilder) SetMaxFileSize(maxFileSize uint64) {
	alb.maxFileSize = maxFileSize
}

func (alb *AccessloggerBuilder) SetMaxFileCount(maxFileCount uint16) {
	alb.maxFileCount = maxFileCount
}

func (alb *AccessloggerBuilder) SetFilename(filename string) {
	alb.filename = filename
}

func (alb *AccessloggerBuilder) SetBufferSize(bufferSize int) {
	alb.bufferSize = bufferSize
}

func (alb *AccessloggerBuilder) buildRemote() writer.WriterDecorator {
	if len(alb.raddr) == 0 {
		return nil
	}

	builder := writer.NewRawUdpWriterDecoratorBuilder()
	builder.SetRemoteAddr(alb.raddr)
	builder.SetRemotePort(alb.rport)
	builder.SetLocalAddr(alb.laddr)
	builder.SetLocalPort(alb.lport)

	w := builder.Build()
	if w == nil {
		return nil
	}

	return w
}

func (alb *AccessloggerBuilder) buildFile() writer.WriterDecorator {
	if len(alb.filename) == 0 {
		return nil
	}

	bufferedfileBuilder := writer.NewFileBufferWriterDecoratorBuilder()
	bufferedfileBuilder.SetBufferSize(alb.bufferSize)
	bufferedfileloger := bufferedfileBuilder.Build()
	if bufferedfileloger == nil {
		return nil
	}

	builder := writer.NewRotateFileWriterDecoratorBuilder()
	builder.SetFilename(alb.filename)
	builder.SetMaxFileSize(alb.maxFileSize)
	builder.SetMaxFileCount(alb.maxFileCount)
	writer := builder.Build()
	if writer == nil {
		return nil
	}

	bufferedfileloger.Wrap(writer)

	return bufferedfileloger
}

func (alb *AccessloggerBuilder) Build() writer.Writer {
	accesslogger := alb.buildRemote()
	filelogger := alb.buildFile()

	if accesslogger != nil {
		accesslogger.Wrap(filelogger)
	} else {
		accesslogger = filelogger
	}

	return accesslogger
}

func NewAccessloggerBuilder() *AccessloggerBuilder {
	return &AccessloggerBuilder{}
}
