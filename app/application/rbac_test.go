package application

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linzhengen/ddd-gin-admin/app/domain/rbac"
)

// mockRBACRepo implements rbac.Repository with configurable return values.
type mockRBACRepo struct {
	rbac.Repository
	rolePolicies []string
	userPolicies []string
}

func (m *mockRBACRepo) ListRolesPolicies(_ context.Context) ([]string, error) {
	return m.rolePolicies, nil
}

func (m *mockRBACRepo) ListUsersPolicies(_ context.Context) ([]string, error) {
	return m.userPolicies, nil
}

func TestRBACAdapter_LoadPolicy_WithCommaSeparatedLines(t *testing.T) {
	const modelConf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) == true && keyMatch2(r.obj, p.obj) == true && regexMatch(r.act, p.act) == true || r.sub == "root"
`
	// The repository returns comma-separated policy lines (the format
	// used by the infrastructure layer). LoadPolicyLine must split
	// these into tokens correctly using a CSV reader.
	mockRepo := &mockRBACRepo{
		rolePolicies: []string{
			"p,test-role-001,/api/v1/test,GET",
			"p,test-role-001,/api/v1/rest/:id,(GET)|(PUT)|(DELETE)",
		},
		userPolicies: []string{
			"g,test-user-001,test-role-001",
		},
	}

	adapter := &rbacAdapter{rbacRepo: mockRepo}

	// Create model from string
	model, err := casbinModel.NewModelFromString(modelConf)
	require.NoError(t, err)

	// Create SyncedEnforcer and initialize with model and adapter,
	// matching the production flow from injector/api/casbin.go
	e, err := casbin.NewSyncedEnforcer()
	require.NoError(t, err)
	err = e.InitWithModelAndAdapter(model, adapter)
	require.NoError(t, err)

	// Verify policies via enforcer's model (InitWithModelAndAdapter copies model internally)
	policies, err := e.GetModel().GetPolicy("p", "p")
	require.NoError(t, err)
	require.Len(t, policies, 2, "should have 2 role policies")
	assert.Contains(t, policies, []string{"test-role-001", "/api/v1/test", "GET"})
	assert.Contains(t, policies, []string{"test-role-001", "/api/v1/rest/:id", "(GET)|(PUT)|(DELETE)"})

	groupingPolicies, err := e.GetModel().GetPolicy("g", "g")
	require.NoError(t, err)
	require.Len(t, groupingPolicies, 1, "should have 1 user-role grouping")
	assert.Contains(t, groupingPolicies, []string{"test-user-001", "test-role-001"})

	t.Run("Root user always allowed", func(t *testing.T) {
		ok, err := e.Enforce("root", "/api/v1/test", "GET")
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("Role user allowed matching path and method", func(t *testing.T) {
		ok, err := e.Enforce("test-user-001", "/api/v1/test", "GET")
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("Role user allowed pattern path", func(t *testing.T) {
		ok, err := e.Enforce("test-user-001", "/api/v1/rest/123", "PUT")
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("Role user denied wrong method", func(t *testing.T) {
		ok, err := e.Enforce("test-user-001", "/api/v1/test", "POST")
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("Unknown user denied", func(t *testing.T) {
		ok, err := e.Enforce("nobody", "/api/v1/test", "GET")
		require.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestRBACAdapter_LoadPolicy_EmptyPolicies(t *testing.T) {
	const modelConf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) == true && keyMatch2(r.obj, p.obj) == true && regexMatch(r.act, p.act) == true || r.sub == "root"
`
	mockRepo := &mockRBACRepo{
		rolePolicies: []string{},
		userPolicies: []string{},
	}

	adapter := &rbacAdapter{rbacRepo: mockRepo}
	model, err := casbinModel.NewModelFromString(modelConf)
	require.NoError(t, err)

	err = adapter.LoadPolicy(model)
	require.NoError(t, err)

	// No policies should be loaded
	policies, err := model.GetPolicy("p", "p")
	require.NoError(t, err)
	assert.Empty(t, policies)
	groupingPolicies, err := model.GetPolicy("g", "g")
	require.NoError(t, err)
	assert.Empty(t, groupingPolicies)
}

func TestRBACAdapter_LoadPolicy_OnlyRolePolicies(t *testing.T) {
	const modelConf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) == true
`
	mockRepo := &mockRBACRepo{
		rolePolicies: []string{
			"p,admin,/api/*,GET",
		},
		userPolicies: []string{},
	}

	adapter := &rbacAdapter{rbacRepo: mockRepo}
	model, err := casbinModel.NewModelFromString(modelConf)
	require.NoError(t, err)

	err = adapter.LoadPolicy(model)
	require.NoError(t, err)

	policies, err := model.GetPolicy("p", "p")
	require.NoError(t, err)
	require.Len(t, policies, 1)
	assert.Equal(t, []string{"admin", "/api/*", "GET"}, policies[0])
}

func TestRBACAdapter_LoadPolicy_PolicyLineWithSpaces(t *testing.T) {
	// Verify that LoadPolicyLine trims leading/trailing spaces correctly
	const modelConf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub)
`
	mockRepo := &mockRBACRepo{
		rolePolicies: []string{
			"p, user1, /data, read",
		},
		userPolicies: []string{
			"g, alice, user1",
		},
	}

	adapter := &rbacAdapter{rbacRepo: mockRepo}
	model, err := casbinModel.NewModelFromString(modelConf)
	require.NoError(t, err)

	err = adapter.LoadPolicy(model)
	require.NoError(t, err)

	// CSV reader trims the space after comma by default (TrimLeadingSpace=true in LoadPolicyLine)
	policies, err := model.GetPolicy("p", "p")
	require.NoError(t, err)
	require.Len(t, policies, 1)
	assert.Equal(t, []string{"user1", "/data", "read"}, policies[0])
}
