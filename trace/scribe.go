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
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/spacemonkeygo/monitor/trace/gen-go/scribe"
	"github.com/spacemonkeygo/monitor/trace/gen-go/zipkin"
	"github.com/spacemonkeygo/spacelog"
)

var (
	logger = spacelog.GetLogger()
)

// ScribeCollector matches the TraceCollector interface, but writes directly
// to a connected Scribe socket.
type ScribeCollector struct {
	transport *thrift.TFramedTransport
	client    *scribe.ScribeClient
}

// NewScribeCollector creates a ScribeCollector. scribe_addr is the address
// of the Scribe endpoint, typically "127.0.0.1:9410"
func NewScribeCollector(scribe_addr string) (*ScribeCollector, error) {
	sa, err := net.ResolveTCPAddr("tcp", scribe_addr)
	if err != nil {
		return nil, err
	}
	transport := thrift.NewTFramedTransport(
		thrift.NewTSocketFromAddrTimeout(sa, 10*time.Second))
	err = transport.Open()
	if err != nil {
		return nil, err
	}

	proto := thrift.NewTBinaryProtocolTransport(transport)
	return &ScribeCollector{
		transport: transport,
		client:    scribe.NewScribeClientProtocol(transport, proto, proto)}, nil
}

// Close closes an existing ScribeCollector
func (s *ScribeCollector) Close() error {
	return s.transport.Close()
}

// CollectSerialized will send a serialized zipkin.Span to the Scribe endpoint
func (c *ScribeCollector) CollectSerialized(serialized []byte) error {
	rc, err := c.client.Log([]*scribe.LogEntry{
		{Category: "zipkin",
			Message: base64.StdEncoding.EncodeToString(serialized)}})
	if err != nil {
		return err
	}
	if rc != scribe.ResultCode_OK {
		return fmt.Errorf("scribe result code not OK: %s", rc)
	}
	return nil
}

// Collect will serialize and send a zipkin.Span to the configured Scribe
// endpoint
func (c *ScribeCollector) Collect(s *zipkin.Span) {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolTransport(t)
	err := s.Write(p)
	if err != nil {
		logger.Errore(err)
	} else {
		logger.Errore(c.CollectSerialized(t.Buffer.Bytes()))
	}
}

var _ TraceCollector = (*ScribeCollector)(nil)
