package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"os"
	"testing"
)

func Test_initialize(t *testing.T) {
	req := getCodeGeneratorRequest()
	// 使用默认选项初始化protogen插件
	plugin, err := protogen.Options{}.New(&req)
	assert.Nil(t, err)

	for _, f := range plugin.Files {
		filename := f.GeneratedFilenamePrefix + "_http.pb.go"
		g := plugin.NewGeneratedFile(filename, f.GoImportPath)
		var fd fileDesc
		fd.initialize(f, g)
		b, _ := json.Marshal(fd)
		fmt.Println("===================", string(b))
	}
}

func Test_generateFile(t *testing.T) {
	req := getCodeGeneratorRequest()
	// 使用默认选项初始化protogen插件
	plugin, err := protogen.Options{}.New(&req)
	assert.Nil(t, err)

	for _, f := range plugin.Files {
		if f.GoPackageName != "v1" {
			continue
		}
		// f.GeneratedFilenamePrefix    github.com/examples/protogen/api/blog/v1/blog
		// f.GoImportPath               github.com/examples/protogen/api/blog/v1
		g := generateFile(plugin, f, true)
		fi, err := os.OpenFile("/tmp/blog_go_http.go", os.O_CREATE|os.O_RDWR, 0644)
		assert.Nil(t, err)
		b, _ := g.Content()
		_, err = fi.Write(b)
		assert.Nil(t, err)
	}
}

func Test_buildPathVars(t *testing.T) {
	newstr := func(s string) *string {
		str := new(string)
		*str = s
		return str
	}
	tests := []struct {
		name    string
		path    string
		wantRes map[string]*string
	}{
		{
			name:    "1",
			path:    "/test/noparams",
			wantRes: map[string]*string{},
		},
		{
			name: "2",
			path: "/test/{message.id}",
			wantRes: map[string]*string{
				"message.id": nil,
			},
		},
		{
			name: "3",
			path: "/test/{message.id}/{message.name=messages/*}",
			wantRes: map[string]*string{
				"message.id":   nil,
				"message.name": newstr("messages/*"),
			},
		},
		{
			name: "4",
			path: "/test/{message.name=messages/*}/books",
			wantRes: map[string]*string{
				"message.name": newstr("messages/*"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPathVars(tt.path)
			assert.Equal(t, tt.wantRes, got)
		})
	}
}

func Test_replacePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantRes string
	}{
		{
			name:    "1",
			path:    "/test/noparams",
			wantRes: "/test/noparams",
		},
		{
			name:    "2",
			path:    "/test/{message.id}",
			wantRes: "/test/{message.id}",
		},
		{
			name:    "3",
			path:    "/test/{message.id=test}",
			wantRes: "/test/{message.id:test}",
		},
		{
			name:    "4",
			path:    "/test/{message.name=messages/*}/books",
			wantRes: "/test/{message.name:messages/.*}/books",
		},
		{
			name:    "5",
			path:    "/test/{message.id}/{message.name=messages/*}",
			wantRes: "/test/{message.id}/{message.name:messages/.*}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := buildPathVars(tt.path)
			for v, s := range vars {
				if s == nil {
					continue
				}
				tt.path = replacePath(v, *s, tt.path)
			}
			assert.Equal(t, tt.wantRes, tt.path)
		})
	}
}

func Test_hasHTTPRule(t *testing.T) {
	f := getProtogenFile("v1")
	tests := []struct {
		name     string
		services []*protogen.Service
		wantRes  bool
	}{
		{
			name:     "1",
			services: f.Services,
			wantRes:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasHTTPRule(tt.services)
			assert.Equal(t, tt.wantRes, got)
		})
	}
}

// Protoc 将protobuf文件编译为 pluginpb.CodeGeneratorRequest结构，并输出到stdin中，
// 这里不通过protoc命令，直接构造
// 根据 api/blog/v1/blog.proto 生成 pluginpb.CodeGeneratorRequest 对象
func getCodeGeneratorRequest() pluginpb.CodeGeneratorRequest {
	b, err := os.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	var req pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(b, &req)
	if err != nil {
		panic(err)
	}
	return req
}

// 根据 api/blog/v1/blog.proto 生成 protogen.File
func getProtogenFile(filename string) *protogen.File {
	req := getCodeGeneratorRequest()
	plugin, err := protogen.Options{}.New(&req)
	if err != nil {
		panic(err)
	}
	for _, f := range plugin.Files {
		if f.GoPackageName == "v1" {
			return f
		}
	}
	return nil
}
