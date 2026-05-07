package userrole

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRoles_ToRoleIDs(t *testing.T) {
	urs := UserRoles{
		{UserID: "u1", RoleID: "r1"},
		{UserID: "u1", RoleID: "r2"},
		{UserID: "u2", RoleID: "r1"},
	}
	ids := urs.ToRoleIDs()
	assert.Equal(t, 3, len(ids))
	assert.Equal(t, []string{"r1", "r2", "r1"}, ids)
}

func TestUserRoles_ToUserIDMap(t *testing.T) {
	urs := UserRoles{
		{UserID: "u1", RoleID: "r1"},
		{UserID: "u1", RoleID: "r2"},
		{UserID: "u2", RoleID: "r1"},
	}
	m := urs.ToUserIDMap()
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 2, len(m["u1"]))
	assert.Equal(t, 1, len(m["u2"]))
}

func TestUserRoles_ToMap_CompositeKey(t *testing.T) {
	urs := UserRoles{
		{UserID: "u1", RoleID: "r1"},
		{UserID: "u1", RoleID: "r2"},
		{UserID: "u2", RoleID: "r1"},
	}
	m := urs.ToMap()
	// Composite key should ensure all entries are preserved
	assert.Equal(t, 3, len(m))
	assert.NotNil(t, m["u1:r1"])
	assert.NotNil(t, m["u1:r2"])
	assert.NotNil(t, m["u2:r1"])
}

func TestUserRoles_ToMap_Empty(t *testing.T) {
	urs := UserRoles{}
	m := urs.ToMap()
	assert.Equal(t, 0, len(m))
}
