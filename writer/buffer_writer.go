package writer

import (
	"bufio"
	"fmt"
	"sync"
)

type writerWrapper struct {
	writer Writer
}

func (ww *writerWrapper) Write(p []byte) (int, error) {
	err := ww.writer.Write("", string(p), "")
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// 使用该Writer会导致后续的Writer的prefix以及suffix消息丢失
type FileBufferWriterDecorator struct {
	buffer *bufio.Writer
	lock   *sync.Mutex
	next   Writer
}

func (ww *FileBufferWriterDecorator) Write(prefix, msg, suffix string) error {
	ww.lock.Lock()
	defer ww.lock.Unlock()

	_, err := ww.buffer.Write([]byte(fmt.Sprint(prefix, msg, suffix)))
	return err
}

func (ww *FileBufferWriterDecorator) Flush() {
	ww.lock.Lock()
	defer ww.lock.Unlock()

	ww.buffer.Flush()
}

func (ww *FileBufferWriterDecorator) Wrap(w Writer) {
	ww.buffer.Flush()
	ww.buffer.Reset(&writerWrapper{writer: w})
	ww.next = w
}

func (ww *FileBufferWriterDecorator) Close() {
	ww.buffer.Flush()
	ww.next.Close()
}

type FileBufferWriterDecoratorBuilder struct {
	bufferSize int
	w          Writer
}

func (wwb *FileBufferWriterDecoratorBuilder) SetBufferSize(bufferSize int) {
	wwb.bufferSize = bufferSize
}

func (wwb *FileBufferWriterDecoratorBuilder) SetNext(w Writer) {
	wwb.w = w
}

func (wwb *FileBufferWriterDecoratorBuilder) Build() *FileBufferWriterDecorator {
	if wwb.bufferSize == 0 {
		return nil
	}

	if wwb.w == nil {
		wwb.w = NewNullWriter()
	}

	buffer := bufio.NewWriterSize(&writerWrapper{writer: &NullWriter{}}, wwb.bufferSize)
	writer := &FileBufferWriterDecorator{
		buffer: buffer,
		lock:   &sync.Mutex{},
	}
	writer.Wrap(wwb.w)

	return writer
}

func NewFileBufferWriterDecoratorBuilder() *FileBufferWriterDecoratorBuilder {
	return &FileBufferWriterDecoratorBuilder{
		w: NewNullWriter(),
	}
}
