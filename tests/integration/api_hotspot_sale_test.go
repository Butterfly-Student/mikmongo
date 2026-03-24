//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
)

// directInsertHotspotSale inserts a hotspot_sale row directly via the repository.
// Returns the sale ID.
func directInsertHotspotSale(t *testing.T, suite *TestSuite, routerID string, agentID *string, profile, batchCode string) string {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	sale := &model.HotspotSale{
		RouterID:     routerID,
		Username:     "v-" + uuid.New().String()[:8],
		Profile:      profile,
		Price:        5000,
		SellingPrice: 7000,
		BatchCode:    batchCode,
		SalesAgentID: agentID,
	}
	err := repos.HotspotSaleRepo.Create(suite.Ctx, sale)
	require.NoError(t, err)
	return sale.ID
}

// hotspotSaleFixture inserts router + agent + 3 sales, returns IDs and admin token.
func hotspotSaleFixture(t *testing.T, suite *TestSuite) (routerID, agentID string, saleIDs []string, adminToken string) {
	t.Helper()

	routerID = createTestMikrotikRouter(t, suite)
	agentStr := createTestSalesAgent(t, suite, routerID)
	agentID = agentStr.ID
	agentIDPtr := &agentID

	saleIDs = []string{
		directInsertHotspotSale(t, suite, routerID, agentIDPtr, "10mb", "BCH1"),
		directInsertHotspotSale(t, suite, routerID, agentIDPtr, "20mb", "BCH1"),
		directInsertHotspotSale(t, suite, routerID, nil, "5mb", "BCH2"),
	}

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	adminToken = loginAs(t, r, email, password)

	return routerID, agentID, saleIDs, adminToken
}

func TestAPIHotspotSale_List_Empty(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	// Use a random router_id that has no sales
	randomID := uuid.New().String()
	w := makeRequest(t, r, http.MethodGet, "/api/v1/hotspot-sales?router_id="+randomID, token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(0), meta["total"])
}

func TestAPIHotspotSale_List(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	routerID, _, _, adminToken := hotspotSaleFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/hotspot-sales?router_id="+routerID, adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(3), meta["total"])
}

func TestAPIHotspotSale_List_FilterAgentID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	_, agentID, _, adminToken := hotspotSaleFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/hotspot-sales?agent_id="+agentID, adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	// 2 sales were created with this agent
	assert.Equal(t, float64(2), meta["total"])
}

func TestAPIHotspotSale_List_FilterProfile(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	routerID, _, _, adminToken := hotspotSaleFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, fmt.Sprintf("/api/v1/hotspot-sales?router_id=%s&profile=10mb", routerID), adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(1), meta["total"])
}

func TestAPIHotspotSale_List_FilterBatchCode(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	routerID, _, _, adminToken := hotspotSaleFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, fmt.Sprintf("/api/v1/hotspot-sales?router_id=%s&batch_code=BCH1", routerID), adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["total"])
}

func TestAPIHotspotSale_List_FilterDate(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	routerID, _, _, adminToken := hotspotSaleFixture(t, suite)
	r := buildTestRouter(t, suite)

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	path := fmt.Sprintf("/api/v1/hotspot-sales?router_id=%s&date_from=%s&date_to=%s", routerID, yesterday, tomorrow)
	w := makeRequest(t, r, http.MethodGet, path, adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(3), meta["total"])
}

func TestAPIHotspotSale_List_InvalidRouterID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/hotspot-sales?router_id=not-a-uuid", token, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIHotspotSale_List_InvalidDate(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/hotspot-sales?date_from=01-13-2024", token, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIHotspotSale_ListByRouter(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	routerID, _, _, adminToken := hotspotSaleFixture(t, suite)
	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/routers/"+routerID+"/hotspot-sales", adminToken, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(3), meta["total"])
}

func TestAPIHotspotSale_ListByRouter_InvalidID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/routers/not-a-uuid/hotspot-sales", token, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIHotspotSale_Unauthorized(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	w := makeRequest(t, r, http.MethodGet, "/api/v1/hotspot-sales", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
