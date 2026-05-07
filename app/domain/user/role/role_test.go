package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoles_ToMap(t *testing.T) {
	roles := Roles{
		{ID: "1", Name: "admin"},
		{ID: "2", Name: "user"},
	}
	m := roles.ToMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, "admin", m["1"].Name)
	assert.Equal(t, "user", m["2"].Name)
}

func TestRoles_ToMap_Empty(t *testing.T) {
	roles := Roles{}
	m := roles.ToMap()
	assert.Equal(t, 0, len(m))
}
