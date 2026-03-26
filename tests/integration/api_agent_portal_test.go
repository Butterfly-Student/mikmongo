//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/model"
)

// agentPortalFixture creates a router, a sales agent via DB, and returns
// the agent ID, username, password, an admin token, and the test engine.
func agentPortalFixture(t *testing.T, suite *TestSuite) (agentID, username, password, adminToken string, engine *gin.Engine) {
	t.Helper()

	engine = buildTestRouter(t, suite)
	email, adminPass, _ := createAPIUser(t, suite, "admin")
	adminToken = loginAs(t, engine, email, adminPass)

	routerID := createTestMikrotikRouter(t, suite)

	username = "portal-" + uuid.New().String()[:8]
	password = "secret123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)

	agent := &model.SalesAgent{
		ID:            uuid.New().String(),
		RouterID:      routerID,
		Name:          "Portal Test Agent",
		Username:      username,
		PasswordHash:  string(hash),
		Status:        "active",
		VoucherMode:   "mix",
		VoucherType:   "upp",
		VoucherLength: 6,
		BillingCycle:  "monthly",
		BillingDay:    1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = suite.DB.WithContext(suite.Ctx).Create(agent).Error
	require.NoError(t, err)

	return agent.ID, username, password, adminToken, engine
}

// agentPortalLoginToken logs in as an agent and returns the access token.
func agentPortalLoginToken(t *testing.T, engine *gin.Engine, username, password string) string {
	t.Helper()
	w := makeRequest(t, engine, http.MethodPost, "/agent-portal/v1/login", "", map[string]string{
		"username": username,
		"password": password,
	})
	require.Equal(t, http.StatusOK, w.Code, "agent login failed: %s", w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)
	token, ok := data["token"].(string)
	require.True(t, ok, "token not found in agent login response")
	return token
}

func TestAgentPortalLogin_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, username, password, _, engine := agentPortalFixture(t, suite)
	w := makeRequest(t, engine, http.MethodPost, "/agent-portal/v1/login", "", map[string]string{
		"username": username,
		"password": password,
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)

	assert.NotEmpty(t, data["token"])
	agent := data["agent"].(map[string]any)
	assert.Equal(t, username, agent["Username"])
	assert.Empty(t, agent["PasswordHash"], "PasswordHash must not be exposed")
}

func TestAgentPortalLogin_InvalidCreds(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, username, _, _, engine := agentPortalFixture(t, suite)
	w := makeRequest(t, engine, http.MethodPost, "/agent-portal/v1/login", "", map[string]string{
		"username": username,
		"password": "wrong-password",
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAgentPortalLogin_InactiveAgent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	agentID, username, _, _, engine := agentPortalFixture(t, suite)

	// Deactivate agent
	err := suite.DB.WithContext(suite.Ctx).
		Model(&model.SalesAgent{}).
		Where("id = ?", agentID).
		Update("status", "inactive").Error
	require.NoError(t, err)

	w := makeRequest(t, engine, http.MethodPost, "/agent-portal/v1/login", "", map[string]string{
		"username": username,
		"password": "secret123",
	})
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAgentPortalGetProfile(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, username, password, _, engine := agentPortalFixture(t, suite)
	token := agentPortalLoginToken(t, engine, username, password)

	w := makeRequest(t, engine, http.MethodGet, "/agent-portal/v1/profile", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)
	assert.Equal(t, username, data["Username"])
	assert.Empty(t, data["PasswordHash"])
}

func TestAgentPortalGetProfile_NoAuth(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, _, _, _, engine := agentPortalFixture(t, suite)

	w := makeRequest(t, engine, http.MethodGet, "/agent-portal/v1/profile", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAgentPortalChangePassword(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, username, password, _, engine := agentPortalFixture(t, suite)
	token := agentPortalLoginToken(t, engine, username, password)

	newPassword := "newpassword456"
	w := makeRequest(t, engine, http.MethodPut, "/agent-portal/v1/profile/password", token, map[string]string{
		"password": newPassword,
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Verify new password works
	w = makeRequest(t, engine, http.MethodPost, "/agent-portal/v1/login", "", map[string]string{
		"username": username,
		"password": newPassword,
	})
	assert.Equal(t, http.StatusOK, w.Code, "login with new password should succeed")
}

func TestAgentPortalGetSales(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	agentID, username, password, _, engine := agentPortalFixture(t, suite)
	token := agentPortalLoginToken(t, engine, username, password)

	// Get agent's router ID from the DB for creating a sale
	var agent model.SalesAgent
	err := suite.DB.WithContext(suite.Ctx).First(&agent, "id = ?", agentID).Error
	require.NoError(t, err)

	createTestHotspotSale(t, suite, agent.RouterID, &agentID, "profile1", "batch1")

	w := makeRequest(t, engine, http.MethodGet, "/agent-portal/v1/sales", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, resp["success"].(bool))
}
