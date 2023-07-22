package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"

	"github.com/meiigo/gkit/log"
)

// Config is a config interface.
type Config interface {
	Load() error
	Scan(v any) error
	Source() ([]byte, error)
}

type config struct {
	opts   options
	values map[string]any
}

// New a config with options.
func New(opts ...Option) Config {
	o := options{
		decoder: defaultDecoder,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &config{
		opts:   o,
		values: make(map[string]any),
	}

}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		kvs, err := src.Load()
		if err != nil {
			return err
		}
		for _, v := range kvs {
			log.Debugf("config loaded: %s format: %s", v.Key, v.Format)
		}
		merged, err := c.Merge(kvs...)
		if err != nil {
			log.Errorf("failed to merge config source: %v", err)
			return err
		}
		c.values = merged
	}
	return nil
}

func (c *config) Scan(v any) error {
	data, err := c.Source()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *config) Source() ([]byte, error) {
	return json.Marshal(c.values)
}

func cloneMap(src map[string]any) (map[string]any, error) {
	var buf bytes.Buffer
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(src)
	if err != nil {
		return nil, err
	}
	var c map[string]any
	err = dec.Decode(&c)
	return c, err
}

func (c *config) Merge(kvs ...*KeyValue) (map[string]any, error) {
	merged, err := cloneMap(c.values)
	if err != nil {
		return merged, err
	}
	for _, kv := range kvs {
		next := make(map[string]any)
		if err := c.opts.decoder(kv, next); err != nil {
			return merged, err
		}
		merged = merge(merged, next)
	}
	return merged, nil
}

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
// merge src into dst
func merge(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = map[string]any{}
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := mapify(srcVal)
			dstMap, dstMapOk := mapify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = merge(dstMap, srcMap)
			}
		}
		dst[key] = srcVal
	}
	return dst
}

func mapify(i any) (map[string]any, bool) {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Map {
		m := make(map[string]any)
		for _, k := range val.MapKeys() {
			m[k.String()] = val.MapIndex(k).Interface()
		}
		return m, true
	}
	return map[string]any{}, false
}
