package config

import (
	"fmt"

	"github.com/meiigo/gkit/x/codec"
)

// Decoder is config decoder.
type Decoder func(*KeyValue, map[string]any) error

// Option is config option.
type Option func(*options)

type options struct {
	sources []Source
	decoder Decoder
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

func WithDecoder(d Decoder) Option {
	return func(o *options) {
		o.decoder = d
	}
}

func defaultDecoder(src *KeyValue, dst map[string]any) error {
	if c := codec.GetCodec(src.Format); c != nil {
		return c.Unmarshal(src.Value, &dst)
	}
	return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
}
