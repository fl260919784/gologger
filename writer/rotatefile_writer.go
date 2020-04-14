package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type RotateFileWriterDecorator struct {
	dirname      string
	filename     string
	maxFileSize  uint64
	maxFileCount uint16
	next         Writer

	lock       *sync.Mutex
	fileWriter *SimpleFileWriter
	fileNames  []string //已经rotate的文件名

}

func (rf *RotateFileWriterDecorator) Write(prefix, msg, suffix string) error {
	if err := rf.fileWriter.Write(prefix, msg, suffix); err != nil {
		return err
	}

	if rf.maxFileSize != 0 {
		rf.check()
	}

	return nil
}

func (rf *RotateFileWriterDecorator) Flush() {
	rf.fileWriter.Flush()
}

func (rf *RotateFileWriterDecorator) Wrap(w Writer) {
	rf.next = w
	rf.fileWriter.Wrap(w)
}

func (rf *RotateFileWriterDecorator) Close() {
	rf.fileWriter.Close()
}

// 识别日志文件是否被删除，是则重建
// 识别文件大小是否超过上限，是则作切割
func (rf *RotateFileWriterDecorator) check() error {
	file, err := os.Stat(rf.filename)
	if err != nil {
		if os.IsNotExist(err) {
			rf.fileWriter.Reopen()
			return nil
		}

		return err
	}

	if uint64(file.Size()) >= rf.maxFileSize {
		return rf.rotate()
	}

	return nil
}

func (rf *RotateFileWriterDecorator) rotate() error {
	rf.lock.Lock()
	defer rf.lock.Unlock()

	now := time.Now()
	suffix := fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	backfilename := fmt.Sprintf("%s.%s", rf.filename, suffix)
	rf.fileNames = append(rf.fileNames, backfilename)
	os.Rename(filepath.Join(rf.dirname, rf.filename), filepath.Join(rf.dirname, backfilename))
	rf.clear()

	rf.fileWriter.Reopen()

	return nil
}

// 按maxFileCount清理过多的历史文件的文件
func (rf *RotateFileWriterDecorator) clear() error {
	size := len(rf.fileNames) - int(rf.maxFileCount)
	if size <= 0 {
		return nil
	}

	for _, filename := range rf.fileNames[0:size] {
		os.Remove(filepath.Join(rf.dirname, filename))
	}

	rf.fileNames = rf.fileNames[size:]

	return nil
}

// 识别当前路径下日志文件，并做相关清理
func (rf *RotateFileWriterDecorator) load() error {
	if rf.maxFileCount == 0 {
		return nil
	}

	rf.lock.Lock()
	defer rf.lock.Unlock()

	oldfilenames := make([]string, 0, rf.maxFileCount)
	files, _ := ioutil.ReadDir(rf.dirname)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), rf.filename) {
			oldfilenames = append(oldfilenames, filepath.Base(file.Name()))
		}
	}
	sort.Sort(sort.StringSlice(oldfilenames))
	rf.fileNames = oldfilenames
	rf.clear()

	return nil
}

type RotateFileWriterDecoratorBuilder struct {
	filename     string // 日志文件名
	dirname      string // 日志文件所在路径名
	maxFileSize  uint64 // 单个文件大小上限
	maxFileCount uint16 // 备份日志文件最大个数
	w            Writer
}

func (rfb *RotateFileWriterDecoratorBuilder) SetFilename(filename string) {
	rfb.dirname = filepath.Dir(filename)
	rfb.filename = filepath.Base(filename)
}

func (rfb *RotateFileWriterDecoratorBuilder) SetMaxFileSize(size uint64) {
	rfb.maxFileSize = size
}

func (rfb *RotateFileWriterDecoratorBuilder) SetMaxFileCount(count uint16) {
	rfb.maxFileCount = count
}

func (rfb *RotateFileWriterDecoratorBuilder) SetNext(w Writer) {
	rfb.w = w
}

func (rfb *RotateFileWriterDecoratorBuilder) Build() *RotateFileWriterDecorator {
	if len(rfb.dirname) == 0 || len(rfb.filename) == 0 {
		return nil
	}

	if rfb.maxFileSize == 0 {
		return nil
	}

	if rfb.w == nil {
		rfb.w = NewNullWriter()
	}

	sfb := SimpleFileWriterBuilder{}
	sfb.SetFilename(filepath.Join(rfb.dirname, rfb.filename))
	sf := sfb.Build()
	if sf == nil {
		return nil
	}

	rf := &RotateFileWriterDecorator{
		dirname:      rfb.dirname,
		filename:     rfb.filename,
		maxFileSize:  rfb.maxFileSize,
		maxFileCount: rfb.maxFileCount,

		lock:       &sync.Mutex{},
		fileWriter: sf,
		fileNames:  make([]string, 0, rfb.maxFileCount),
	}

	rf.Wrap(rfb.w)
	rf.load()

	return rf
}

func NewRotateFileWriterDecoratorBuilder() *RotateFileWriterDecoratorBuilder {
	return &RotateFileWriterDecoratorBuilder{
		w: NewNullWriter(),
	}
}
