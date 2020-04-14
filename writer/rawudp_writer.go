package writer

import (
	"fmt"
	"net"
)

type RawUdpWriterDecorator struct {
	laddr string // local addr
	raddr string // remote addr
	conn  *net.UDPConn
	next  Writer
}

func (ruw *RawUdpWriterDecorator) Write(prefix, msg, suffix string) error {
	if _, err := ruw.conn.Write([]byte(msg)); err != nil {
		return err
	}

	return ruw.next.Write(prefix, msg, suffix)
}

func (ruw *RawUdpWriterDecorator) Flush() {
	ruw.next.Flush()
}

func (ruw *RawUdpWriterDecorator) Wrap(w Writer) {
	ruw.next = w
}

func (ruw *RawUdpWriterDecorator) Close() {
	ruw.conn.Close()
	ruw.next.Close()
}

type RawUdpWriterDecoratorBuilder struct {
	raddr string
	rport int

	laddr string
	lport int
	w     Writer
}

func (ruwb *RawUdpWriterDecoratorBuilder) SetRemoteAddr(addr string) {
	ruwb.raddr = addr
}

func (ruwb *RawUdpWriterDecoratorBuilder) SetRemotePort(port int) {
	ruwb.rport = port
}

func (ruwb *RawUdpWriterDecoratorBuilder) SetLocalAddr(addr string) {
	ruwb.laddr = addr
}

func (ruwb *RawUdpWriterDecoratorBuilder) SetLocalPort(port int) {
	ruwb.lport = port
}

func (ruwb *RawUdpWriterDecoratorBuilder) SetNext(w Writer) {
	ruwb.w = w
}

func (ruwb *RawUdpWriterDecoratorBuilder) Build() *RawUdpWriterDecorator {
	raddstr := net.JoinHostPort(ruwb.raddr, fmt.Sprint(ruwb.rport))
	laddstr := net.JoinHostPort(ruwb.raddr, fmt.Sprint(ruwb.rport))

	raddr, err := net.ResolveUDPAddr("udp", raddstr)
	if err != nil {
		return nil
	}

	laddr, err := net.ResolveUDPAddr("udp", laddstr)
	if err != nil {
		return nil
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return nil
	}

	if ruwb.w == nil {
		wwb.w = NewNullWriter()
	}

	writer := &RawUdpWriterDecorator{
		laddr: raddstr,
		raddr: laddstr,
		conn:  conn,
	}
	writer.Wrap(wwb.w)

	return writer
}

func NewRawUdpWriterDecoratorBuilder() *RawUdpWriterDecoratorBuilder {
	return &RawUdpWriterDecoratorBuilder{
		w: NewNullWriter(),
	}
}
