# 概要说明



实现基础的日志库，根据业务场景不同，分离出不同的日志后端

通过不同的后端的灵活组合产生更多的特性



Writer为基础接口

WriterDecorator接口目的让各种后端可随意组合

Logger为日志外壳



# 种类说明



| 类名                | 接口            | 功能                     |
| ------------------- | --------------- | ------------------------ |
| NullWriter          | Writer          | 丢弃所有写入             |
| SimpleFileWriterDecorator    | WriterDecorator | 普通文件方式             |
| AutoCheckFileWriterDecorator | WriterDecorator | 自动检测日志文件是否存在 |
| RotateFileWriterDecorator    | WriterDecorator | 自动切割日志文件方式     |
| FileBufferWriterDecorator    | WriterDecorator | buffer方式               |
| RawUdpWriterDecorator        | WriterDecorator | udp方式                  |



## 备注

FileBufferWriterDecorator因为Buffer的存在，所以会合并perfix、message、suffix，这对启Wrap的对象有要求

即Wrap对象的message能忍受此行为，如RawUdpWriterDecorator则不适合作为其Wrap对象，因为RawUdpWriter只需要Message



# 组合类

## AccessloggerBuilder

组合RawUdpWriter、FileBufferWriter、RotateFileWriter



## BufferedFileLogerBuilder

组合FileBufferWriter、RotateFileWriter



# 开箱即用

default_logger.go

通过实现全局默认的Logger对象

内部Writer为SimpleFileWriter，直接输出到标准输出，用户可重新设置

并开放Debug、Infof、Warnf、Errorf对外提供开箱即用



default_accessslogger.go

类似提供全局默认的Logger对象

内部Writer为SimpleFileWriter，直接输出到标准输出，用户可重新设置

