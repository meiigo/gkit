package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_defaultDecoder(t *testing.T) {
	tests := []struct {
		name string
		src  *KeyValue
		dst  map[string]any
	}{
		{
			name: "1",
			src: &KeyValue{
				Key:    "yaml",
				Value:  []byte(_testYaml),
				Format: "yaml",
			},
			dst: map[string]any{
				"Foo": map[string]any{
					"bar": []map[string]any{
						{
							"name": "nihao",
							"age":  1,
						},
						{
							"name": "nihao",
							"age":  1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := defaultDecoder(tt.src, tt.dst)
			assert.Nil(t, err)
		})
	}
}
