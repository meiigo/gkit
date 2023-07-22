package grpc

import (
	"github.com/meiigo/gkit/transport"
	"google.golang.org/grpc/metadata"
)

var _ transport.Transporter = &Transport{}

type Transport struct {
	endpoint    string
	operation   string
	reqHeader   grpcHeader
	replyHeader grpcHeader
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.KindGRPC
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

type grpcHeader metadata.MD

// Get returns the value associated with the passed key.
func (hd grpcHeader) Get(key string) string {
	vals := metadata.MD(hd).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (hd grpcHeader) Set(key string, value string) {
	metadata.MD(hd).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hd grpcHeader) Keys() []string {
	keys := make([]string, 0, len(hd))
	for k := range metadata.MD(hd) {
		keys = append(keys, k)
	}
	return keys
}
