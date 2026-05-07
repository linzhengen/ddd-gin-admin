package response

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domainrole "github.com/linzhengen/ddd-gin-admin/app/domain/user/role"
)

func TestRolesFromDomain_NilPointerRegression(t *testing.T) {
	roles := domainrole.Roles{
		{ID: "1", Name: "admin"},
		{ID: "2", Name: "user"},
	}
	result := RolesFromDomain(roles)
	assert.Equal(t, 2, len(result))
	assert.NotNil(t, result[0])
	assert.NotNil(t, result[1])
	assert.Equal(t, "admin", result[0].Name)
	assert.Equal(t, "user", result[1].Name)
}

func TestRolesFromDomain_EmptyInput(t *testing.T) {
	result := RolesFromDomain(domainrole.Roles{})
	assert.Equal(t, 0, len(result))
}

func TestRolesFromDomain_NilInput(t *testing.T) {
	result := RolesFromDomain(nil)
	assert.Equal(t, 0, len(result))
}
