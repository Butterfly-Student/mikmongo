//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// salesAgentFixture creates a router, one sales agent via API, and an admin user.
// Returns routerID, agentID, and admin token.
func salesAgentFixture(t *testing.T, suite *TestSuite) (routerID, agentID, adminToken string) {
	t.Helper()

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	adminToken = loginAs(t, r, email, password)

	routerID = createTestMikrotikRouter(t, suite)

	body := map[string]interface{}{
		"router_id":      routerID,
		"name":           "API Agent",
		"username":       "apiagent-" + uuid.New().String()[:8],
		"password":       "secret123",
		"status":         "active",
		"voucher_mode":   "mix",
		"voucher_type":   "upp",
		"voucher_length": 6,
	}
	w := makeRequest(t, r, http.MethodPost, "/api/v1/sales-agents", adminToken, body)
	require.Equal(t, http.StatusCreated, w.Code, "fixture agent creation failed: %s", w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	agentID = data["ID"].(string) // PascalCase — model has no json tags

	return routerID, agentID, adminToken
}

func TestAPICreateSalesAgent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)
	routerID := createTestMikrotikRouter(t, suite)

	body := map[string]interface{}{
		"router_id":      routerID,
		"name":           "Test Agent",
		"username":       "testagent01",
		"password":       "secure123",
		"status":         "active",
		"voucher_mode":   "mix",
		"voucher_type":   "upp",
		"voucher_length": 8,
		"bill_discount":  5.0,
	}
	w := makeRequest(t, r, http.MethodPost, "/api/v1/sales-agents", token, body)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})

	assert.NotEmpty(t, data["ID"])
	assert.Equal(t, "Test Agent", data["Name"])
	assert.Equal(t, "testagent01", data["Username"])
	assert.Empty(t, data["PasswordHash"], "PasswordHash must not be exposed")
}

func TestAPICreateSalesAgent_ShortPassword(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)
	routerID := createTestMikrotikRouter(t, suite)

	body := map[string]interface{}{
		"router_id": routerID,
		"name":      "Agent",
		"username":  "agentoops",
		"password":  "abc", // too short (min=6)
	}
	w := makeRequest(t, r, http.MethodPost, "/api/v1/sales-agents", token, body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPICreateSalesAgent_MissingRequired(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/sales-agents", token, map[string]interface{}{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIGetSalesAgent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/sales-agents/"+agentID, adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})

	assert.Equal(t, agentID, data["ID"])
	assert.Empty(t, data["PasswordHash"], "PasswordHash must not be exposed")
}

func TestAPIGetSalesAgent_NotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, _, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/sales-agents/"+uuid.New().String(), adminToken, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIGetSalesAgent_InvalidID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, _, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/sales-agents/not-a-uuid", adminToken, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIListSalesAgents(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	routerID, _, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, fmt.Sprintf("/api/v1/sales-agents?router_id=%s", routerID), adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	meta := resp["meta"].(map[string]interface{})
	assert.GreaterOrEqual(t, int(meta["total"].(float64)), 1)
}

func TestAPIListSalesAgents_NoFilter(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, _, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/sales-agents", adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestAPIUpdateSalesAgent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	newName := "Updated Agent Name"
	w := makeRequest(t, r, http.MethodPut, "/api/v1/sales-agents/"+agentID, adminToken, map[string]interface{}{
		"name":   newName,
		"status": "inactive",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, newName, data["Name"])
	assert.Equal(t, "inactive", data["Status"])
	assert.Empty(t, data["PasswordHash"])
}

func TestAPIDeleteSalesAgent(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodDelete, "/api/v1/sales-agents/"+agentID, adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestAPIDeleteSalesAgent_GetAfterDelete(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	makeRequest(t, r, http.MethodDelete, "/api/v1/sales-agents/"+agentID, adminToken, nil)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/sales-agents/"+agentID, adminToken, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIUpsertProfilePrice_Create(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	body := map[string]interface{}{
		"base_price":    5000,
		"selling_price": 7000,
		"is_active":     true,
	}
	path := fmt.Sprintf("/api/v1/sales-agents/%s/profile-prices/10mb", agentID)
	w := makeRequest(t, r, http.MethodPut, path, adminToken, body)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, 5000.0, data["BasePrice"])
	assert.Equal(t, 7000.0, data["SellingPrice"])
}

func TestAPIUpsertProfilePrice_Update(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	path := fmt.Sprintf("/api/v1/sales-agents/%s/profile-prices/10mb", agentID)

	makeRequest(t, r, http.MethodPut, path, adminToken, map[string]interface{}{
		"base_price": 5000, "selling_price": 7000, "is_active": true,
	})

	w := makeRequest(t, r, http.MethodPut, path, adminToken, map[string]interface{}{
		"base_price": 8000, "selling_price": 10000, "is_active": true,
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, 8000.0, data["BasePrice"])
	assert.Equal(t, 10000.0, data["SellingPrice"])
}

func TestAPIListProfilePrices(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, adminToken := salesAgentFixture(t, suite)
	r := buildTestRouter(t, suite)

	for _, profile := range []string{"10mb", "20mb"} {
		path := fmt.Sprintf("/api/v1/sales-agents/%s/profile-prices/%s", agentID, profile)
		makeRequest(t, r, http.MethodPut, path, adminToken, map[string]interface{}{
			"base_price": 5000, "selling_price": 7000, "is_active": true,
		})
	}

	w := makeRequest(t, r, http.MethodGet, fmt.Sprintf("/api/v1/sales-agents/%s/profile-prices", agentID), adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestAPISalesAgent_Unauthorized(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	w := makeRequest(t, r, http.MethodGet, "/api/v1/sales-agents", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
