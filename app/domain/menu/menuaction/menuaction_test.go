package menuaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenuActions_ToMenuIDMap(t *testing.T) {
	actions := MenuActions{
		{ID: "a1", MenuID: "m1", Code: "read"},
		{ID: "a2", MenuID: "m1", Code: "write"},
		{ID: "a3", MenuID: "m2", Code: "read"},
	}
	m := actions.ToMenuIDMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 2, len(m["m1"]))
	assert.Equal(t, 1, len(m["m2"]))
}

func TestMenuActions_ToMap(t *testing.T) {
	actions := MenuActions{
		{Code: "read", Name: "Read"},
		{Code: "write", Name: "Write"},
	}
	m := actions.ToMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, "Read", m["read"].Name)
	assert.Equal(t, "Write", m["write"].Name)
}
