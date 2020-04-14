package writer

type Writer interface {
	Write(prefix, msg, suffix string) error
	Flush()
	Close()
}

type WriterDecorator interface {
	Writer
	Wrap(Writer)
}

type NullWriter struct {
}

func (nw NullWriter) Write(prefix, msg, suffix string) error {
	return nil
}

func (nw NullWriter) Flush() {
}

func (nw NullWriter) Close() {
}

func NewNullWriter() *NullWriter {
	return &NullWriter{}
}

type Builder interface {
	Build() Writer
}
