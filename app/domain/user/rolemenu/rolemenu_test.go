package rolemenu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleMenus_ToMap(t *testing.T) {
	rms := RoleMenus{
		{RoleID: "r1", MenuID: "m1", ActionID: "a1"},
		{RoleID: "r1", MenuID: "m1", ActionID: "a2"},
	}
	m := rms.ToMap()
	assert.Equal(t, 2, len(m))
	assert.NotNil(t, m["m1-a1"])
	assert.NotNil(t, m["m1-a2"])
}

func TestRoleMenus_ToMenuIDs_Unique(t *testing.T) {
	rms := RoleMenus{
		{MenuID: "m1", ActionID: "a1"},
		{MenuID: "m1", ActionID: "a2"},
		{MenuID: "m2", ActionID: "a3"},
	}
	ids := rms.ToMenuIDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, "m1")
	assert.Contains(t, ids, "m2")
}

func TestRoleMenus_ToActionIDs_Unique(t *testing.T) {
	rms := RoleMenus{
		{ActionID: "a1"},
		{ActionID: "a1"},
		{ActionID: "a2"},
	}
	ids := rms.ToActionIDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, "a1")
	assert.Contains(t, ids, "a2")
}

func TestRoleMenus_ToRoleIDMap(t *testing.T) {
	rms := RoleMenus{
		{RoleID: "r1", MenuID: "m1"},
		{RoleID: "r1", MenuID: "m2"},
		{RoleID: "r2", MenuID: "m3"},
	}
	m := rms.ToRoleIDMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 2, len(m["r1"]))
	assert.Equal(t, 1, len(m["r2"]))
}
