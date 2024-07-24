package types_test

import (
	"testing"

	"github.com/korrel8r/client/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	q := types.ParseQuery("x:y:z")
	assert.Equal(t, "x:y:z", q.String())
	assert.Equal(t, []string{"x", "y", "z"}, []string{q.Class.Domain, q.Class.Name, q.Data})
}

func TestClass(t *testing.T) {
	c := types.ParseClass("x:y")
	assert.Equal(t, "x:y", c.String())
	assert.Equal(t, c.Domain, "x")
	assert.Equal(t, c.Name, "y")
}
