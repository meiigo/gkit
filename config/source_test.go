package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testJSON = `
{
    "server":{
        "http":{
            "addr":"0.0.0.0",
			"port":80,
            "timeout":0.5,
			"enable_ssl":true
        },
        "grpc":{
            "addr":"0.0.0.0",
			"port":10080,
            "timeout":0.2
        }
    },
    "data":{
        "database":{
            "driver":"mysql",
            "source":"root:root@tcp(127.0.0.1:3306)/migo?parseTime=true"
        }
    },
	"endpoints":[
		"www.aaa.com",
		"www.bbb.org"
	]
}`

	_testFooJson = `{
  "Foo": {
    "bar": [
      {
        "age": 11,
        "name": "world"
      },
      {
        "age": 22
      }
    ]
  }
}`

	_testYaml = `
Foo:
  bar:
    - {name: hello,age: 1}
    - {name: hello,age: 2}
`
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		source  Source
		wantKey string
		wantVal []byte
	}{
		{
			name:    "json",
			source:  newTestJsonSource(_testFooJson),
			wantKey: "json",
			wantVal: []byte(_testFooJson),
		},
		{
			name:    "yaml",
			source:  newTestYamlSource(_testYaml),
			wantKey: "yaml",
			wantVal: []byte(_testYaml),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.source.Load()
			assert.Nil(t, err)
			assert.Equal(t, tt.wantKey, got[0].Key)
			assert.Equal(t, tt.wantVal, got[0].Value)
		})
	}
}

type testJSONSource struct {
	data string
}

func newTestJsonSource(data string) *testJSONSource {
	return &testJSONSource{data: data}
}

func (p *testJSONSource) Load() ([]*KeyValue, error) {
	kv := &KeyValue{
		Key:    "json",
		Value:  []byte(p.data),
		Format: "json",
	}
	return []*KeyValue{kv}, nil
}

type testYamlSource struct {
	data string
}

func newTestYamlSource(data string) *testYamlSource {
	return &testYamlSource{data: data}
}

func (p *testYamlSource) Load() ([]*KeyValue, error) {
	kv := &KeyValue{
		Key:    "yaml",
		Value:  []byte(p.data),
		Format: "yaml",
	}
	return []*KeyValue{kv}, nil
}
