package main

import (
	"bytes"
	"html/template"
	"strings"
)

var httpTemplate = `
{{$svrType := .ServiceType}}

// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v1.0.0
// source: {{.SourceFilePath}}

package {{.PackageName}}

import (
    "context"
    "github.com/gin-gonic/gin"
)

var _ = new(context.Context)

type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}HTTPHandler(g *gin.Engine, srv {{.ServiceType}}HTTPServer) {
	{{- range .Methods}}
	g.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) func(c *gin.Context) {
	return func(c *gin.Context) {
        ctx := c.Request.Context()
		var in {{.Request}}
		if err := c.ShouldBind(&in); err != nil {
            c.AbortWithStatusJSON(400, gin.H{"err": err.Error()})
            return
		}

        // execute
        out, err := srv.{{.Name}}(ctx, &in)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
            return
        }
		c.JSON(200, out)
	}
}
{{end}}
`

type fileDesc struct {
	Version        string // v1.0.0
	SourceFilePath string // api/hello/hello.proto
	PackageName    string // v1
	serviceDesc
}

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // hello.Greeter
	// Metadata    string // api/hello/hello.proto
	Methods    []*methodDesc
	MethodSets map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name         string
	OriginalName string // The parsed original name
	Num          int
	Request      string
	Reply        string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (f *fileDesc) execute() string {
	f.MethodSets = make(map[string]*methodDesc)
	for _, m := range f.Methods {
		f.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	//tmpl, err := template.ParseFiles(tpl)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
