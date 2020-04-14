package writer

import (
	"os"
)

type AutoCheckFileWriterDecorator struct {
	filename   string
	fileWriter *SimpleFileWriterDecorator
	next       Writer
}

func (ac *AutoCheckFileWriterDecorator) Write(prefix, msg, suffix string) error {
	err := ac.fileWriter.Write(prefix, msg, suffix)
	if err != nil {
		return err
	}

	ac.check()
	return nil
}

func (ac *AutoCheckFileWriterDecorator) Flush() {
	ac.fileWriter.Flush()
}

func (ac *AutoCheckFileWriterDecorator) Wrap(w Writer) {
	ac.next = w
	ac.fileWriter.Wrap(w)
}

func (ac *AutoCheckFileWriterDecorator) Close() {
	ac.fileWriter.Close()
}

// 识别日志文件是否存在
// 不存在则创建
func (ac *AutoCheckFileWriterDecorator) check() error {
	if _, err := os.Stat(ac.filename); os.IsNotExist(err) {
		return ac.fileWriter.Reopen()
	}

	return nil
}

// AutoCheckFileWriterDecorator的创建类
// 主要用于收集目标类的属性，然后创建
type AutoCheckFileWriterDecoratorBuilder struct {
	filename string
	w        Writer
}

func (acfb *AutoCheckFileWriterDecoratorBuilder) SetFilename(filename string) {
	acfb.filename = filename
}

func (acfb *AutoCheckFileWriterDecoratorBuilder) SetNext(w Writer) {
	acfb.w = w
}

func (acfb *AutoCheckFileWriterDecoratorBuilder) Build() *AutoCheckFileWriterDecorator {
	if len(acfb.filename) == 0 {
		return nil
	}

	if acfb.w == nil {
		acfb.w = NewNullWriter()
	}

	sfb := NewSimpleFileWriterDecoratorBuilder()
	sfb.SetFilename(acfb.filename)
	sf := sfb.Build()
	if sf == nil {
		return nil
	}

	ac := &AutoCheckFileWriterDecorator{
		filename:   acfb.filename,
		fileWriter: sf,
	}

	ac.Wrap(acfb.w)

	if ac.check() != nil {
		return nil
	}

	return ac
}

func NewAutoCheckFileWriterDecoratorBuilder() *AutoCheckFileWriterDecoratorBuilder {
	return &AutoCheckFileWriterDecoratorBuilder{
		w: NewNullWriter(),
	}
}
