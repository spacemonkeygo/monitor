// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"net"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/spacemonkeygo/monitor/trace/gen-go/zipkin"
)

const (
	maxPacketSize = 8192
)

// UDPCollector matches the TraceCollector interface, but sends serialized
// zipkin.Span objects over UDP, instead of the Scribe protocol. See
// RedirectPackets for the UDP server-side code.
type UDPCollector struct {
	ch   chan *zipkin.Span
	conn *net.UDPConn
	addr *net.UDPAddr
}

// NewUDPCollector creates a UDPCollector that sends packets to collector_addr.
// buffer_size is how many outstanding unsent zipkin.Span objects can exist
// before Spans start getting dropped.
func NewUDPCollector(collector_addr string, buffer_size int) (
	*UDPCollector, error) {
	addr, err := net.ResolveUDPAddr("udp", collector_addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	c := &UDPCollector{
		ch:   make(chan *zipkin.Span, buffer_size),
		conn: conn,
		addr: addr}
	go c.handle()
	return c, nil
}

func (c *UDPCollector) handle() {
	for {
		select {
		case s, ok := <-c.ch:
			if !ok {
				return
			}
			logger.Errore(c.send(s))
		}
	}
}

func (c *UDPCollector) send(s *zipkin.Span) error {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolTransport(t)
	err := s.Write(p)
	if err != nil {
		return err
	}
	_, err = c.conn.WriteToUDP(t.Buffer.Bytes(), c.addr)
	return err
}

// Collect takes a zipkin.Span object, serializes it, and sends it to the
// configured collector_addr.
func (c *UDPCollector) Collect(span *zipkin.Span) {
	select {
	case c.ch <- span:
	default:
	}
}

// RedirectPackets is a method that handles incoming packets from the
// UDPCollector class. RedirectPackets, when running, will listen for UDP
// packets containing serialized zipkin.Span objects on listen_addr, then will
// resend those packets to the given ScribeCollector. On any error,
// RedirectPackets currently aborts.
func RedirectPackets(listen_addr string, collector *ScribeCollector) error {
	la, err := net.ResolveUDPAddr("udp", listen_addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", la)
	if err != nil {
		return err
	}
	defer conn.Close()
	var buf [maxPacketSize]byte
	for {
		n, _, err := conn.ReadFrom(buf[:])
		if err != nil {
			return err
		}
		err = collector.CollectSerialized(buf[:n])
		if err != nil {
			return err
		}
	}
}
