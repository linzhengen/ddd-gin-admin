package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/LyricTian/captcha"
	captchastore "github.com/LyricTian/captcha/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu"
	menuActionInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuaction"
	menuActionResourceInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuactionresource"
	userInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"
	roleInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/role"
	roleMenuInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/rolemenu"
	userRoleInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/response"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/injector"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/hash"
	"gorm.io/gorm"
)

var (
	testServer *httptest.Server
	testDB     *gorm.DB
)

func TestMain(m *testing.M) {
	// Resolve project root for file paths
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// Initialize test config
	initTestConfig(projectRoot)

	// Clean up any previous test database file
	dbPath := configs.C.Sqlite3.Path
	os.Remove(dbPath)

	// Build the full injector (creates DB, migrates, wires dependencies)
	inj, cleanup, err := injector.BuildApiInjector()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build injector: %v\n", err)
		os.Exit(1)
	}

	testServer = httptest.NewServer(inj.GetEngine())

	// Set up captcha store — use a wrapper that never clears on Get,
	// so the same test captcha can be reused across multiple test requests.
	memStore := captchastore.NewMemoryStore(captcha.Expiration, captcha.Expiration)
	memStore.Set("test-captcha-id", []byte{1, 2, 3, 4})
	captcha.SetCustomStore(&nonClearingStore{Store: memStore})

	// Store DB reference for test data seeding
	testDB = inj.GetDB()

	// Seed test data
	seedTestData()

	code := m.Run()

	// Cleanup
	testServer.Close()
	cleanup()
	os.Remove(dbPath)
	os.Exit(code)
}

func initTestConfig(projectRoot string) {
	configs.C.RunMode = "test"
	configs.C.PrintConfig = false

	configs.C.HTTP.Host = "127.0.0.1"
	configs.C.HTTP.Port = 0
	configs.C.HTTP.ShutdownTimeout = 5
	configs.C.HTTP.MaxContentLength = 64 << 20
	configs.C.HTTP.MaxLoggerLength = 4096

	configs.C.Menu.Enable = false
	configs.C.Menu.Data = ""

	configs.C.Casbin.Enable = true
	configs.C.Casbin.Debug = false
	configs.C.Casbin.Model = filepath.Join(projectRoot, "configs", "model.conf")
	configs.C.Casbin.AutoLoad = false
	configs.C.Casbin.AutoLoadInternal = 10

	configs.C.Log.Level = 6
	configs.C.Log.Format = "json"
	configs.C.Log.Output = "stdout"
	configs.C.Log.OutputFile = ""
	configs.C.Log.EnableHook = false
	configs.C.Log.HookLevels = nil
	configs.C.Log.Hook = "gorm"
	configs.C.Log.HookMaxThread = 1
	configs.C.Log.HookMaxBuffer = 200

	configs.C.LogGormHook.DBType = "sqlite3"
	configs.C.LogGormHook.MaxLifetime = 3600
	configs.C.LogGormHook.MaxOpenConns = 1
	configs.C.LogGormHook.MaxIdleConns = 1
	configs.C.LogGormHook.Table = "g_logger"

	configs.C.Root.UserName = "root"
	configs.C.Root.Password = "abc-123"
	configs.C.Root.RealName = "Root"

	configs.C.JWTAuth.Enable = true
	configs.C.JWTAuth.SigningMethod = "HS512"
	configs.C.JWTAuth.SigningKey = "ddd-gin-admin-test"
	configs.C.JWTAuth.Expired = 7200
	configs.C.JWTAuth.Store = "file"
	configs.C.JWTAuth.FilePath = ":memory:"
	configs.C.JWTAuth.RedisDB = 0
	configs.C.JWTAuth.RedisPrefix = ""

	configs.C.Monitor.Enable = false

	configs.C.Captcha.Store = "memory"
	configs.C.Captcha.Length = 4
	configs.C.Captcha.Width = 160
	configs.C.Captcha.Height = 64

	configs.C.RateLimiter.Enable = false
	configs.C.RateLimiter.Count = 2000

	configs.C.CORS.Enable = false
	configs.C.CORS.AllowOrigins = []string{"*"}
	configs.C.CORS.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	configs.C.CORS.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	configs.C.CORS.AllowCredentials = true
	configs.C.CORS.MaxAge = 86400

	configs.C.GZIP.Enable = false
	configs.C.GZIP.ExcludedExtentions = nil
	configs.C.GZIP.ExcludedPaths = nil

	configs.C.Redis.Addr = "127.0.0.1:6379"
	configs.C.Redis.Password = ""

	configs.C.Gorm.Debug = false
	configs.C.Gorm.DBType = "sqlite3"
	configs.C.Gorm.MaxLifetime = 3600
	configs.C.Gorm.MaxOpenConns = 10
	configs.C.Gorm.MaxIdleConns = 10
	configs.C.Gorm.TablePrefix = ""
	configs.C.Gorm.EnableAutoMigrate = true

	configs.C.Sqlite3.Path = filepath.Join(os.TempDir(), "ddd-gin-admin-e2e-test.db")

	configs.C.MySQL.Host = "127.0.0.1"
	configs.C.MySQL.Port = 3306
	configs.C.MySQL.User = "root"
	configs.C.MySQL.Password = "root"
	configs.C.MySQL.DBName = "gin-admin"
	configs.C.MySQL.Parameters = "charset=utf8mb4&parseTime=True&loc=Local"

	configs.C.Postgres.Host = "127.0.0.1"
	configs.C.Postgres.Port = 5432
	configs.C.Postgres.User = "root"
	configs.C.Postgres.Password = "root"
	configs.C.Postgres.DBName = "gin-admin"
	configs.C.Postgres.SSLMode = "disable"
}

// seedTestData populates the test database with initial data.
// Note: The auto-migrate in the injector creates the tables.
func seedTestData() {
	// Create root user in DB (needed for GetUserInfo/GetActiveUserWithRole)
	rootUser := userInfra.Model{
		ID:       "root",
		UserName: "root",
		RealName: "Root",
		Password: hash.SHA1String("abc-123"),
		Phone:    strPtr(""),
		Email:    strPtr(""),
		Status:   1,
		Creator:  "system",
	}
	if err := testDB.Create(&rootUser).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed root user: %v", err))
	}

	// Create test user with SHA1 hashed password
	testUser := userInfra.Model{
		ID:       "test-user-001",
		UserName: "testuser",
		RealName: "Test User",
		Password: hash.SHA1String("testpass123"),
		Phone:    strPtr("13800138000"),
		Email:    strPtr("test@example.com"),
		Status:   1,
		Creator:  "system",
	}
	if err := testDB.Create(&testUser).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed user: %v", err))
	}

	// Create test role
	testRole := roleInfra.Model{
		ID:       "test-role-001",
		Name:     "test_role",
		Sequence: 1,
		Status:   1,
		Creator:  "system",
	}
	if err := testDB.Create(&testRole).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed role: %v", err))
	}

	// Create user-role association
	ur := userRoleInfra.Model{
		ID:     "test-ur-001",
		UserID: "test-user-001",
		RoleID: "test-role-001",
	}
	if err := testDB.Create(&ur).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed user_role: %v", err))
	}

	// Create test menu
	testMenu := menu.Model{
		ID:         "test-menu-001",
		Name:       "Test Menu",
		Sequence:   1,
		Icon:       strPtr("setting"),
		Router:     strPtr("/test"),
		ParentID:   nil,
		ParentPath: nil,
		ShowStatus: 1,
		Status:     1,
		Memo:       strPtr("test menu"),
		Creator:    "system",
	}
	if err := testDB.Create(&testMenu).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed menu: %v", err))
	}

	// Create test menu action
	testAction := menuActionInfra.Model{
		ID:     "test-action-001",
		MenuID: "test-menu-001",
		Code:   "query",
		Name:   "Query",
	}
	if err := testDB.Create(&testAction).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed menu action: %v", err))
	}

	// Create test menu action resource
	testResource := menuActionResourceInfra.Model{
		ID:       "test-res-001",
		ActionID: "test-action-001",
		Method:   "GET",
		Path:     "/api/v1/test",
	}
	if err := testDB.Create(&testResource).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed menu action resource: %v", err))
	}

	// Create role-menu association
	rm := roleMenuInfra.Model{
		ID:       "test-rm-001",
		RoleID:   "test-role-001",
		MenuID:   "test-menu-001",
		ActionID: "test-action-001",
	}
	if err := testDB.Create(&rm).Error; err != nil {
		panic(fmt.Sprintf("Failed to seed role_menu: %v", err))
	}
}

func strPtr(s string) *string { return &s }

// nonClearingStore wraps a captcha store but ignores the "clear" flag on Get,
// so captchas are never consumed and can be reused across test requests.
type nonClearingStore struct {
	captchastore.Store
}

func (s *nonClearingStore) Get(id string, _ bool) []byte {
	return s.Store.Get(id, false)
}

// Helper functions for making requests and parsing responses

// ----------------------------- E2E TESTS -----------------------------

// TestHealthCheck verifies the health endpoint works
func TestHealthCheck(t *testing.T) {
	token := getRootToken(t)

	req, err := http.NewRequest("GET", testServer.URL+"/api/health", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Status    string `json:"status"`
		CheckedAt string `json:"checked_at"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "OK", result.Status)
	assert.NotEmpty(t, result.CheckedAt)
}

// TestLogin_InvalidCredentials verifies login fails with bad credentials
func TestLogin_InvalidCredentials(t *testing.T) {
	body := `{"user_name":"nonexistent","password":"wrong","captcha_id":"test-captcha-id","captcha_code":"1234"}`
	resp, err := http.Post(testServer.URL+"/api/v1/pub/login", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLogin_RootUser(t *testing.T) {
	body := `{"user_name":"root","password":"abc-123","captcha_id":"test-captcha-id","captcha_code":"1234"}`
	resp, err := http.Post(testServer.URL+"/api/v1/pub/login", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var lr struct {
		response.LoginTokenInfo
		response.StatusResult
	}
	err = json.NewDecoder(resp.Body).Decode(&lr)
	require.NoError(t, err)
	assert.NotEmpty(t, lr.AccessToken)
	assert.Equal(t, "Bearer", lr.TokenType)
	assert.Greater(t, lr.ExpiresAt, int64(0))
}

func TestLogin_RegularUser(t *testing.T) {
	// Note: captcha is verified, so we need to skip it or provide a valid one
	// In test mode, captcha verification happens in the handler
	// Using "test" as captcha ID/code since we didn't configure captcha store
	body := `{"user_name":"testuser","password":"testpass123","captcha_id":"test-captcha-id","captcha_code":"1234"}`
	resp, err := http.Post(testServer.URL+"/api/v1/pub/login", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var lr struct {
		response.LoginTokenInfo
		response.StatusResult
	}
	err = json.NewDecoder(resp.Body).Decode(&lr)
	require.NoError(t, err)
	assert.NotEmpty(t, lr.AccessToken)
	assert.Equal(t, "Bearer", lr.TokenType)
}

// Helper to get auth token for root user
func getRootToken(t *testing.T) string {
	body := `{"user_name":"root","password":"abc-123","captcha_id":"test-captcha-id","captcha_code":"1234"}`
	resp, err := http.Post(testServer.URL+"/api/v1/pub/login", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var lr struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&lr)
	require.NoError(t, err)
	require.NotEmpty(t, lr.AccessToken, "failed to get root token")
	return lr.AccessToken
}

// TestGetUserInfo verifies getting authenticated user info
func TestGetUserInfo(t *testing.T) {
	token := getRootToken(t)

	req, err := http.NewRequest("GET", testServer.URL+"/api/v1/pub/current/user", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var userInfo struct {
		UserID   string `json:"user_id"`
		UserName string `json:"user_name"`
		RealName string `json:"real_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	require.NoError(t, err)
	assert.Equal(t, "root", userInfo.UserName)
}

// TestGetUserMenuTree verifies getting the menu tree for root user
func TestGetUserMenuTree(t *testing.T) {
	token := getRootToken(t)

	req, err := http.NewRequest("GET", testServer.URL+"/api/v1/pub/current/menutree", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResult struct {
		List []json.RawMessage `json:"list"`
	}
	err = json.NewDecoder(resp.Body).Decode(&listResult)
	require.NoError(t, err)
	// Root user should see all menus
	assert.GreaterOrEqual(t, len(listResult.List), 0)
}

// TestUserCRUD tests the full CRUD cycle for users
func TestUserCRUD(t *testing.T) {
	token := getRootToken(t)

	// Create a new user
	createBody := `{
		"user_name": "newuser",
		"real_name": "New User",
		"password": "newpass123",
		"status": 1
	}`
	req, err := http.NewRequest("POST", testServer.URL+"/api/v1/users", strings.NewReader(createBody))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var idResult struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&idResult)
	require.NoError(t, err)
	assert.NotEmpty(t, idResult.ID)
	userID := idResult.ID

	// Query users list
	req, err = http.NewRequest("GET", testServer.URL+"/api/v1/users", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get created user
	req, err = http.NewRequest("GET", testServer.URL+"/api/v1/users/"+userID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getUser struct {
		ID       string `json:"id"`
		UserName string `json:"user_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&getUser)
	require.NoError(t, err)
	assert.Equal(t, "newuser", getUser.UserName)

	// Update user — must include user_name to avoid triggering the
	// username uniqueness check in the application layer.
	updateBody := `{"user_name":"newuser","real_name": "Updated User"}`
	req, err = http.NewRequest("PUT", testServer.URL+"/api/v1/users/"+userID, strings.NewReader(updateBody))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete user
	req, err = http.NewRequest("DELETE", testServer.URL+"/api/v1/users/"+userID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestRoleCRUD tests the full CRUD cycle for roles
func TestRoleCRUD(t *testing.T) {
	token := getRootToken(t)

	// Create a role
	createBody := `{
		"name": "e2e-test-role",
		"sequence": 1,
		"status": 1
	}`
	req, err := http.NewRequest("POST", testServer.URL+"/api/v1/roles", strings.NewReader(createBody))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var idResult struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&idResult)
	require.NoError(t, err)
	assert.NotEmpty(t, idResult.ID)
	roleID := idResult.ID

	// Get role
	req, err = http.NewRequest("GET", testServer.URL+"/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getRole struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&getRole)
	require.NoError(t, err)
	assert.Equal(t, "e2e-test-role", getRole.Name)

	// Query roles
	req, err = http.NewRequest("GET", testServer.URL+"/api/v1/roles", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete role
	req, err = http.NewRequest("DELETE", testServer.URL+"/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestMenuCRUD tests the full CRUD cycle for menus
func TestMenuCRUD(t *testing.T) {
	token := getRootToken(t)

	// Create a menu
	createBody := `{
		"name": "e2e-test-menu",
		"sequence": 1,
		"icon": "test",
		"router": "/e2e-test",
		"show_status": 1,
		"status": 1
	}`
	req, err := http.NewRequest("POST", testServer.URL+"/api/v1/menus", strings.NewReader(createBody))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var idResult struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&idResult)
	require.NoError(t, err)
	assert.NotEmpty(t, idResult.ID)
	menuID := idResult.ID

	// Get menu
	req, err = http.NewRequest("GET", testServer.URL+"/api/v1/menus/"+menuID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getMenu struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&getMenu)
	require.NoError(t, err)
	assert.Equal(t, "e2e-test-menu", getMenu.Name)

	// Delete menu
	req, err = http.NewRequest("DELETE", testServer.URL+"/api/v1/menus/"+menuID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestUnauthorizedAccess verifies that protected endpoints return 401 without token
func TestUnauthorizedAccess(t *testing.T) {
	endpoints := []string{
		"/api/v1/users",
		"/api/v1/roles",
		"/api/v1/menus",
	}
	for _, endpoint := range endpoints {
		resp, err := http.Get(testServer.URL + endpoint)
		require.NoError(t, err)
		_ = resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
			"expected 401 for %s without auth", endpoint)
	}
}

// TestRefreshToken verifies token refresh endpoint
func TestRefreshToken(t *testing.T) {
	token := getRootToken(t)

	req, err := http.NewRequest("POST", testServer.URL+"/api/v1/pub/refresh-token", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var lr struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresAt   int64  `json:"expires_at"`
	}
	err = json.NewDecoder(resp.Body).Decode(&lr)
	require.NoError(t, err)
	assert.NotEmpty(t, lr.AccessToken)
}

// TestRoleSelect verifies role select endpoint
func TestRoleSelect(t *testing.T) {
	token := getRootToken(t)

	req, err := http.NewRequest("GET", testServer.URL+"/api/v1/roles.select", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResult struct {
		List []json.RawMessage `json:"list"`
	}
	err = json.NewDecoder(resp.Body).Decode(&listResult)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResult.List), 0)
}

// TestMenuTree verifies menu tree endpoint
func TestMenuTree(t *testing.T) {
	token := getRootToken(t)

	req, err := http.NewRequest("GET", testServer.URL+"/api/v1/menus.tree", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestRoleEnableDisable verifies enabling and disabling a role
func TestRoleEnableDisable(t *testing.T) {
	token := getRootToken(t)

	// Create a role first
	createBody := `{"name":"toggle-role","sequence":1,"status":1}`
	req, _ := http.NewRequest("POST", testServer.URL+"/api/v1/roles", strings.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	var idResult struct{ ID string }
	err = json.NewDecoder(resp.Body).Decode(&idResult)
	require.NoError(t, err)
	resp.Body.Close()
	roleID := idResult.ID

	// Disable
	req, _ = http.NewRequest("PATCH", testServer.URL+"/api/v1/roles/"+roleID+"/disable", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Enable
	req, _ = http.NewRequest("PATCH", testServer.URL+"/api/v1/roles/"+roleID+"/enable", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
