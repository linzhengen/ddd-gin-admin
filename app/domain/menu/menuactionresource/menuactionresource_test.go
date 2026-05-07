package menuactionresource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenuActionResources_ToMap_KeyCollision(t *testing.T) {
	resources := MenuActionResources{
		&MenuActionResource{Method: "GET", Path: "/api/users"},
		&MenuActionResource{Method: "POST", Path: "/api/users"},
		&MenuActionResource{Method: "GET", Path: "/api/roles"},
	}
	m := resources.ToMap()
	// With separator, "GET|/api/users" and "POST|/api/users" are distinct keys
	assert.Equal(t, 3, len(m))
	assert.NotNil(t, m["GET|/api/users"])
	assert.NotNil(t, m["POST|/api/users"])
	assert.NotNil(t, m["GET|/api/roles"])
}

func TestMenuActionResources_ToActionIDMap(t *testing.T) {
	resources := MenuActionResources{
		&MenuActionResource{ID: "1", ActionID: "action1", Method: "GET", Path: "/api/users"},
		&MenuActionResource{ID: "2", ActionID: "action1", Method: "POST", Path: "/api/users"},
		&MenuActionResource{ID: "3", ActionID: "action2", Method: "GET", Path: "/api/roles"},
	}
	m := resources.ToActionIDMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 2, len(m["action1"]))
	assert.Equal(t, 1, len(m["action2"]))
}

func TestMenuActionResources_ToMenuActionIDMap(t *testing.T) {
	resources := MenuActionResources{
		&MenuActionResource{ActionID: "a1"},
		&MenuActionResource{ActionID: "a1"},
		&MenuActionResource{ActionID: "a2"},
	}
	m := resources.ToMenuActionIDMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 2, len(m["a1"]))
	assert.Equal(t, 1, len(m["a2"]))
}
