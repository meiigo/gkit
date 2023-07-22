package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	o := &options{}
	v := "123"
	ID(v)(o)
	assert.Equal(t, v, o.id)

}
