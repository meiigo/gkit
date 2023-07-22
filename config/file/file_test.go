package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meiigo/gkit/config"
)

const (
	_testJSON = `
{
    "test":{
        "settings":{
            "int_key":1000,
            "float_key":1000.1,
            "duration_key":10000,
            "string_key":"string_value"
        },
        "server":{
            "addr":"127.0.0.1",
            "port":8000
        }
    },
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ]
}`
)

// go test -v *.go -test.run=^TestLoadFile$
func TestLoadFile(t *testing.T) {
	var (
		path     = filepath.Join(os.TempDir(), "test_config")
		filename = filepath.Join(path, "test.json")
		data     = []byte(_testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		t.Error(err)
	}

	testSource(t, path, data)
	testSource(t, path, data)
}

func testSource(t *testing.T, path string, data []byte) {
	t.Log(path)
	s := NewSource(path)
	kvs, err := s.Load()
	assert.Nil(t, err)
	assert.Equal(t, kvs[0].Value, data)
}

func TestConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config.json")
	defer os.Remove(path)
	if err := os.WriteFile(path, []byte(_testJSON), 0o666); err != nil {
		t.Error(err)
	}
	c := config.New(config.WithSource(
		NewSource(path),
	))
	testScan(t, c)
}

func testScan(t *testing.T, c config.Config) {
	type TestJSON struct {
		Test struct {
			Settings struct {
				IntKey      int     `json:"int_key"`
				FloatKey    float64 `json:"float_key"`
				DurationKey int     `json:"duration_key"`
				StringKey   string  `json:"string_key"`
			} `json:"settings"`
			Server struct {
				Addr string `json:"addr"`
				Port int    `json:"port"`
			} `json:"server"`
		} `json:"test"`
		Foo []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		} `json:"foo"`
	}
	var conf TestJSON
	err := c.Load()
	assert.Nil(t, err)
	err = c.Scan(&conf)
	assert.Nil(t, err)
}
