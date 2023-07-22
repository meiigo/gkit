package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_serviceDesc_execute(t *testing.T) {
	tests := []struct {
		name string
		desc fileDesc
		want string
	}{
		{
			name: "blog",
			desc: fileDesc{
				SourceFilePath: "api/hello/hello.proto",
				PackageName:    "v1",
				serviceDesc: serviceDesc{
					ServiceType: "BlogService",
					ServiceName: "BlogService",
					Methods: []*methodDesc{
						{
							Name:    "CreateArticle",
							Request: "CreateArticleRequest",
							Reply:   "CreateArticleReply",
							Path:    "/v1/article/",
							Method:  "POST",
						},
						{
							Name:    "GetArticle",
							Request: "GetArticleRequest",
							Reply:   "GetArticleReply",
							Path:    "/v1/article/{id}",
							Method:  "GET",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileDesc{
				SourceFilePath: tt.desc.SourceFilePath,
				PackageName:    tt.desc.PackageName,
				serviceDesc:    tt.desc.serviceDesc,
			}
			got := f.execute()
			fi, err := os.OpenFile("/tmp/tmp.go", os.O_CREATE|os.O_RDWR, 0644)
			assert.Nil(t, err)
			_, err = fi.WriteString(got)
			assert.Nil(t, err)
		})
	}
}
