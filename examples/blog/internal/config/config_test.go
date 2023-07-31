package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "1",
			path: "../../config/local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitConfig(tt.path)
			assert.Nil(t, err)
		})
	}
}
