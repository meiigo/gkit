package config

import (
	"encoding/json"
	"testing"

	_ "github.com/meiigo/gkit/x/codec/json"
	_ "github.com/meiigo/gkit/x/codec/yaml"
	"github.com/stretchr/testify/assert"
)

func Test_config_Load(t *testing.T) {
	tests := []struct {
		name       string
		sources    []Source
		decoder    Decoder
		wantSource string
	}{
		{
			name: "1",
			sources: []Source{
				newTestJsonSource(_testFooJson),
				newTestYamlSource(_testYaml),
			},
			decoder:    defaultDecoder,
			wantSource: `{"Foo":{"bar":[{"age":1,"name":"hello"},{"age":2,"name":"hello"}]}}`,
		},
		{
			name: "2",
			sources: []Source{
				newTestYamlSource(_testYaml),
				newTestJsonSource(_testFooJson),
			},
			decoder:    defaultDecoder,
			wantSource: `{"Foo":{"bar":[{"age":12,"name":"nihao"},{"age":17}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := New(
				WithSource(tt.sources...),
				WithDecoder(tt.decoder),
			)
			err := cf.Load()
			assert.Nil(t, err)
			d, err := cf.Source()
			assert.Nil(t, err)
			assert.Equal(t, tt.wantSource, string(d))
		})
	}
}

func Test_cloneMap(t *testing.T) {
	tests := []struct {
		original     map[string]any
		transformer  func(m map[string]any) map[string]any
		wantCopy     map[string]any
		wantOriginal map[string]any
	}{
		{
			original: nil,
			transformer: func(m map[string]any) map[string]any {
				return map[string]any{}
			},
			wantCopy:     map[string]any{},
			wantOriginal: nil,
		},
		{
			original: map[string]any{},
			transformer: func(m map[string]any) map[string]any {
				return nil
			},
			wantCopy:     nil,
			wantOriginal: map[string]any{},
		},
		{
			original: map[string]any{},
			transformer: func(m map[string]any) map[string]any {
				m["foo"] = "bar"
				return m
			},
			wantCopy:     map[string]any{"foo": "bar"},
			wantOriginal: map[string]any{},
		},
		{
			original: map[string]any{"foo": "bar"},
			transformer: func(m map[string]any) map[string]any {
				m["foo"] = "car"
				return m
			},
			wantCopy:     map[string]any{"foo": "car"},
			wantOriginal: map[string]any{"foo": "bar"},
		},
		{
			original: map[string]any{},
			transformer: func(m map[string]any) map[string]any {
				m["foo"] = map[string]any{"biz": "bar"}
				return m
			},
			wantCopy:     map[string]any{"foo": map[string]any{"biz": "bar"}},
			wantOriginal: map[string]any{},
		},
		{
			original: map[string]any{"foo": []string{"biz", "baz"}},
			transformer: func(m map[string]any) map[string]any {
				m["foo"].([]string)[0] = "hiz"
				return m
			},
			wantCopy:     map[string]any{"foo": []string{"hiz", "baz"}},
			wantOriginal: map[string]any{"foo": []string{"biz", "baz"}},
		},
	}
	for _, tt := range tests {
		got, err := cloneMap(tt.original)
		assert.Nil(t, err)
		assert.Exactly(t, tt.wantCopy, tt.transformer(got))
		assert.Exactly(t, tt.wantOriginal, tt.original)
	}
}

func Test_merge(t *testing.T) {
	tests := []struct {
		dst, src, want map[string]any
	}{
		{
			dst:  nil,
			src:  nil,
			want: map[string]any{},
		},
		{
			dst:  nil,
			src:  map[string]any{},
			want: map[string]any{},
		},
		{
			dst:  map[string]any{},
			src:  nil,
			want: map[string]any{},
		},
		{
			dst:  map[string]any{},
			src:  map[string]any{},
			want: map[string]any{},
		},
		{
			dst: map[string]any{
				"b": "bb",
				"c": 99,
			},
			src: map[string]any{
				"a": 1,
				"b": "b",
			},
			want: map[string]any{
				"a": 1,
				"b": "b",
				"c": 99,
			},
		},
		{
			dst: nil,
			src: map[string]any{
				"a": 1,
				"b": "b",
			},
			want: map[string]any{
				"a": 1,
				"b": "b",
			},
		},
		{
			dst: map[string]any{},
			src: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 1,
					"d": []string{"1", "2", "3"},
					"e": map[string]any{
						"f": 1,
					},
				},
			},
			want: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 1,
					"d": []string{"1", "2", "3"},
					"e": map[string]any{
						"f": 1,
					},
				},
			},
		},
		{
			dst: map[string]any{
				"b": map[string]any{},
			},
			src: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 1,
					"d": []string{"1", "2", "3"},
					"e": map[string]any{
						"f": 1,
					},
				},
			},
			want: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 1,
					"d": []string{"1", "2", "3"},
					"e": map[string]any{
						"f": 1,
					},
				},
			},
		},
		{
			dst: map[string]any{
				"b": map[string]any{
					"e": map[string]any{
						"f": 99,
					},
				},
			},
			src: map[string]any{
				"b": map[string]any{
					"e": map[string]any{
						"f": "99",
					},
				},
			},
			want: map[string]any{
				"b": map[string]any{
					"e": map[string]any{
						"f": "99",
					},
				},
			},
		},
		{
			dst: map[string]any{
				"b": map[string]any{
					"d": []string{"4", "5"},
					"e": map[string]any{
						"f": 99,
					},
				},
			},
			src: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 1,
					"d": []string{"1", "2", "3"},
					"e": map[string]any{
						"f": 1,
						"h": "a",
					},
				},
			},
			want: map[string]any{
				"a": 1,
				"b": map[string]any{
					"c": 1,
					"d": []string{"1", "2", "3"},
					"e": map[string]any{
						"f": 1,
						"h": "a",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got := merge(tt.dst, tt.src)
		assert.Equal(t, tt.want, got)
	}
}

func Test_config_Merge(t *testing.T) {
	tests := []struct {
		name     string
		opts     options
		kvs      []*KeyValue
		wantData string
	}{
		{
			name: "1",
			opts: options{
				decoder: defaultDecoder,
			},
			kvs: []*KeyValue{
				{
					Key:    "yaml",
					Value:  []byte(_testYaml),
					Format: "yaml",
				},
				{
					Key:    "json",
					Value:  []byte(_testFooJson),
					Format: "json",
				},
			},
			wantData: `{"Foo":{"bar":[{"age":12,"name":"nihao"},{"age":17}]}}`,
		},
		{
			name: "2",
			opts: options{
				decoder: defaultDecoder,
			},
			kvs: []*KeyValue{
				{
					Key:    "json",
					Value:  []byte(_testFooJson),
					Format: "json",
				},
				{
					Key:    "yaml",
					Value:  []byte(_testYaml),
					Format: "yaml",
				},
			},
			wantData: `{"Foo":{"bar":[{"age":1,"name":"hello"},{"age":2,"name":"hello"}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := config{
				opts:   tt.opts,
				values: map[string]any{},
			}
			got, err := c.Merge(tt.kvs...)
			assert.Nil(t, err)
			b, _ := json.Marshal(got)
			assert.Equal(t, tt.wantData, string(b))
		})
	}
}

func Test_config_Scan(t *testing.T) {
	type Bar struct {
		Age  int    `json:"age"`
		Name string `json:"name,omitempty"`
	}
	type Foo struct {
		Bars []Bar `json:"bar"`
	}
	type T struct {
		Foo Foo `json:"Foo"`
	}

	tests := []struct {
		name     string
		sources  []Source
		decoder  Decoder
		T, wantT *T
	}{
		{
			name: "1",
			sources: []Source{
				newTestJsonSource(_testFooJson),
				newTestYamlSource(_testYaml),
			},
			decoder: defaultDecoder,
			T:       &T{},
			wantT: &T{
				Foo: Foo{
					Bars: []Bar{
						{Name: "hello", Age: 1},
						{Name: "hello", Age: 2},
					},
				},
			},
		},
		{
			name: "2",
			sources: []Source{
				newTestYamlSource(_testYaml),
				newTestJsonSource(_testFooJson),
			},
			decoder: defaultDecoder,
			T:       &T{},
			wantT: &T{
				Foo: Foo{
					Bars: []Bar{
						{Name: "world", Age: 11},
						{Name: "", Age: 22},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := New(
				WithSource(tt.sources...),
				WithDecoder(tt.decoder),
			)
			err := cf.Load()
			assert.Nil(t, err)
			err = cf.Scan(tt.T)
			assert.Nil(t, err)
			assert.Equal(t, tt.wantT, tt.T)
		})
	}
}
