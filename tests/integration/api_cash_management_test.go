//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cashFixture creates an admin user, logs in, and returns the engine + token + userID.
func cashFixture(t *testing.T, suite *TestSuite) (engine *gin.Engine, token, userID string) {
	t.Helper()
	engine = buildTestRouter(t, suite)
	email, password, userID := createAPIUser(t, suite, "admin")
	token = loginAs(t, engine, email, password)
	return engine, token, userID
}

func TestCashEntryCreateAndGet(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	body := map[string]any{
		"type":           "expense",
		"source":         "operational",
		"amount":         500000,
		"description":    "Sewa tower Maret",
		"payment_method": "bank_transfer",
		"bank_name":      "BCA",
		"entry_date":     "2026-03-15",
	}
	w := makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, body)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)
	entryID := data["ID"].(string)
	assert.NotEmpty(t, entryID)
	assert.Equal(t, "pending", data["Status"])
	assert.Contains(t, data["EntryNumber"], "KAS")

	// GET single entry
	w = makeRequest(t, engine, http.MethodGet, "/api/v1/cash-entries/"+entryID, token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestCashEntryList(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	// Create two entries
	for i := 0; i < 2; i++ {
		body := map[string]any{
			"type":           "income",
			"source":         "other",
			"amount":         100000,
			"description":    fmt.Sprintf("Test income %d", i),
			"payment_method": "cash",
		}
		w := makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, body)
		require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	}

	w := makeRequest(t, engine, http.MethodGet, "/api/v1/cash-entries?type=income", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]any)
	assert.GreaterOrEqual(t, meta["total"].(float64), float64(2))
}

func TestCashEntryApproveAndReject(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	// Create an entry to approve
	body := map[string]any{
		"type":           "expense",
		"source":         "operational",
		"amount":         100000,
		"description":    "Test approve",
		"payment_method": "cash",
	}
	w := makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, body)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	approveID := resp["data"].(map[string]any)["ID"].(string)

	// Approve
	w = makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries/"+approveID+"/approve", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "approved", resp["data"].(map[string]any)["Status"])

	// Create another entry to reject
	body["description"] = "Test reject"
	w = makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, body)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	rejectID := resp["data"].(map[string]any)["ID"].(string)

	// Reject
	w = makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries/"+rejectID+"/reject", token, map[string]string{
		"reason": "not approved by management",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "rejected", resp["data"].(map[string]any)["Status"])
}

func TestCashEntryUpdateAndDelete(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	body := map[string]any{
		"type":           "expense",
		"source":         "operational",
		"amount":         200000,
		"description":    "Original description",
		"payment_method": "cash",
	}
	w := makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, body)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	entryID := resp["data"].(map[string]any)["ID"].(string)

	// Update
	w = makeRequest(t, engine, http.MethodPut, "/api/v1/cash-entries/"+entryID, token, map[string]any{
		"description": "Updated description",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Delete
	w = makeRequest(t, engine, http.MethodDelete, "/api/v1/cash-entries/"+entryID, token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Verify deleted (should be 404)
	w = makeRequest(t, engine, http.MethodGet, "/api/v1/cash-entries/"+entryID, token, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPettyCashFundCRUD(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, userID := cashFixture(t, suite)

	// Create fund
	body := map[string]any{
		"fund_name":       "Kas Kecil Operasional",
		"initial_balance": 5000000,
		"custodian_id":    userID,
	}
	w := makeRequest(t, engine, http.MethodPost, "/api/v1/petty-cash", token, body)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)
	fundID := data["ID"].(string)
	assert.Equal(t, float64(5000000), data["CurrentBalance"])

	// Get fund
	w = makeRequest(t, engine, http.MethodGet, "/api/v1/petty-cash/"+fundID, token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// List funds
	w = makeRequest(t, engine, http.MethodGet, "/api/v1/petty-cash", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	meta := resp["meta"].(map[string]any)
	assert.GreaterOrEqual(t, meta["total"].(float64), float64(1))

	// Update fund
	w = makeRequest(t, engine, http.MethodPut, "/api/v1/petty-cash/"+fundID, token, map[string]any{
		"fund_name": "Kas Kecil Updated",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Top up fund
	w = makeRequest(t, engine, http.MethodPost, "/api/v1/petty-cash/"+fundID+"/topup", token, map[string]any{
		"amount": 1000000,
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data = resp["data"].(map[string]any)
	assert.Equal(t, float64(6000000), data["CurrentBalance"])
}

func TestPettyCashExpenseDebit(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, userID := cashFixture(t, suite)

	// Create fund
	fundBody := map[string]any{
		"fund_name":       "Kas Kecil Test",
		"initial_balance": 1000000,
		"custodian_id":    userID,
	}
	w := makeRequest(t, engine, http.MethodPost, "/api/v1/petty-cash", token, fundBody)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	fundID := resp["data"].(map[string]any)["ID"].(string)

	// Create expense linked to petty cash
	entryBody := map[string]any{
		"type":               "expense",
		"source":             "operational",
		"amount":             300000,
		"description":        "Beli ATK",
		"payment_method":     "cash",
		"petty_cash_fund_id": fundID,
	}
	w = makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, entryBody)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	// Check fund balance was debited
	w = makeRequest(t, engine, http.MethodGet, "/api/v1/petty-cash/"+fundID, token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(700000), resp["data"].(map[string]any)["CurrentBalance"])
}

func TestCashFlowReport(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	// Create an income and expense entry (approved)
	incomeBody := map[string]any{
		"type":           "income",
		"source":         "other",
		"amount":         1000000,
		"description":    "Test income",
		"payment_method": "cash",
		"entry_date":     "2026-03-15",
	}
	w := makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries", token, incomeBody)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	incomeID := resp["data"].(map[string]any)["ID"].(string)

	// Approve income
	w = makeRequest(t, engine, http.MethodPost, "/api/v1/cash-entries/"+incomeID+"/approve", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	// Get cash flow report
	w = makeRequest(t, engine, http.MethodGet, "/api/v1/reports/cash-flow?from=2026-03-01&to=2026-03-31", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)
	assert.GreaterOrEqual(t, data["total_income"].(float64), float64(1000000))
}

func TestCashBalanceReport(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	w := makeRequest(t, engine, http.MethodGet, "/api/v1/reports/cash-balance", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]any)
	assert.NotNil(t, data["balance"])
}

func TestCashEntryNotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	engine, token, _ := cashFixture(t, suite)

	w := makeRequest(t, engine, http.MethodGet, "/api/v1/cash-entries/"+uuid.New().String(), token, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
