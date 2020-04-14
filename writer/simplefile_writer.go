package writer

import (
	"fmt"
	"os"
)

type SimpleFileWriterDecorator struct {
	filename string
	file     *os.File
	next     Writer
}

func (fw *SimpleFileWriterDecorator) Write(prefix, msg, suffix string) error {
	_, err := fw.file.Write([]byte(fmt.Sprint(prefix, msg, suffix)))
	if err != nil {
		return err
	}

	return fw.next.Write(prefix, msg, suffix)
}

func (fw *SimpleFileWriterDecorator) Flush() {
	fw.next.Flush()
}

func (fw *SimpleFileWriterDecorator) Wrap(w Writer) {
	if w == nil {
		return
	}

	fw.next = w
}

func (fw *SimpleFileWriterDecorator) Close() {
	fw.file.Close()
	fw.next.Close()
}

// 重新打开日志文件，若不存在的创建
// 关闭之前的文件句柄
func (fw *SimpleFileWriterDecorator) Reopen() error {
	file, err := os.OpenFile(fw.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if fw.file != nil {
		fw.file.Close()
	}

	fw.file = file
	return nil
}

// SimpleFileWriterDecorator的创建类
// 主要用于收集目标类的属性，然后创建
type SimpleFileWriterDecoratorBuilder struct {
	filename string
	w        Writer
}

func (sfb *SimpleFileWriterDecoratorBuilder) SetFilename(filename string) {
	sfb.filename = filename
}

func (sfb *SimpleFileWriterDecoratorBuilder) SetNext(w Writer) {
	sfb.w = w
}

func (sfb *SimpleFileWriterDecoratorBuilder) Build() *SimpleFileWriterDecorator {
	if len(sfb.filename) == 0 {
		return nil
	}

	if sfb.w == nil {
		sfb.w = NewNullWriter()
	}

	fw := &SimpleFileWriterDecorator{
		filename: sfb.filename,
	}

	fw.Wrap(sfb.w)
	if fw.Reopen() != nil {
		return nil
	}

	return fw
}

func NewSimpleFileWriterDecoratorBuilder() *SimpleFileWriterDecoratorBuilder {
	return &SimpleFileWriterDecoratorBuilder{
		w: NewNullWriter(),
	}
}
