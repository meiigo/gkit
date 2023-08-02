package grpc

import (
	"context"
	"testing"

	"github.com/meiigo/gkit/examples/blog/api"
	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		opts []ClientOption
	}{
		{
			name: "1",
			ctx:  context.Background(),
			opts: []ClientOption{
				WithEndpoint("127.0.0.1:9090"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DialInsecure(tt.ctx, tt.opts...)
			assert.Nil(t, err)
			cli := api.NewBlogClient(got)
			resp, err := cli.GetArticle(tt.ctx, &api.GetArticleRequest{})
			t.Log("resp:", resp)
			t.Log("err:", err)
		})
	}
}
