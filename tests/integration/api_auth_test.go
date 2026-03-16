//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIAuth_Login_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")

	w := makeRequest(t, r, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, data["access_token"])
}

func TestAPIAuth_Login_WrongPassword(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, _, _ := createAPIUser(t, suite, "admin")

	w := makeRequest(t, r, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": "wrongpassword",
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPIAuth_Login_EmptyBody(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)

	// Send malformed JSON (a plain string instead of an object) to trigger ShouldBindJSON error → 400
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`"not_an_object"`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIAuth_GetMe_Authenticated(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/auth/me", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, email, data["Email"])
}

func TestAPIAuth_GetMe_NoToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/auth/me", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPIAuth_GetMe_InvalidToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)

	w := makeRequest(t, r, http.MethodGet, "/api/v1/auth/me", "garbage", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPIAuth_Logout_BlacklistsToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	// Logout
	w := makeRequest(t, r, http.MethodPost, "/api/v1/auth/logout", token, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// Access /me with same token — must be rejected
	w2 := makeRequest(t, r, http.MethodGet, "/api/v1/auth/me", token, nil)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestAPIAuth_Refresh_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	_, refreshToken := loginAsWithRefresh(t, r, email, password)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/auth/refresh", "", map[string]string{
		"refresh_token": refreshToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, data["access_token"])
}

func TestAPIAuth_ChangePassword_Success(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	newPassword := "NewPassword456!"
	w := makeRequest(t, r, http.MethodPost, "/api/v1/auth/change-password", token, map[string]string{
		"old_password": password,
		"new_password": newPassword,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	// Old credentials should fail
	w2 := makeRequest(t, r, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	assert.Equal(t, http.StatusUnauthorized, w2.Code)

	// New credentials should succeed
	w3 := makeRequest(t, r, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    email,
		"password": newPassword,
	})
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestAPIAuth_ChangePassword_WrongOld(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	r := buildTestRouter(t, suite)
	email, password, _ := createAPIUser(t, suite, "admin")
	token := loginAs(t, r, email, password)

	w := makeRequest(t, r, http.MethodPost, "/api/v1/auth/change-password", token, map[string]string{
		"old_password": "wrongoldpassword",
		"new_password": "NewPassword456!",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
