package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func buildRouterIDTestRouter(middleware gin.HandlerFunc, handler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.GET("/test/:router_id", middleware, handler)
	return r
}

func TestValidateRouterID_InvalidUUID(t *testing.T) {
	mw := NewMikrotikRouterMiddleware()

	r := buildRouterIDTestRouter(mw.ValidateRouterID(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid router ID format")
}

func TestValidateRouterID_ValidUUID(t *testing.T) {
	mw := NewMikrotikRouterMiddleware()
	validUUID := uuid.New()

	var gotRouterID interface{}
	handler := func(c *gin.Context) {
		gotRouterID, _ = c.Get("router_id")
		c.Status(http.StatusOK)
	}

	r := buildRouterIDTestRouter(mw.ValidateRouterID(), handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test/"+validUUID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, validUUID, gotRouterID)
}

func TestValidateRouterID_EmptyString(t *testing.T) {
	mw := NewMikrotikRouterMiddleware()

	// With gin, a route parameter /:router_id will never be empty because
	// the route simply won't match. Use a route without the param to test
	// the empty-string path.
	r := gin.New()
	r.GET("/test", mw.ValidateRouterID(), okHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Router ID is required")
}

func TestValidateRouterID_NumericString(t *testing.T) {
	mw := NewMikrotikRouterMiddleware()

	r := buildRouterIDTestRouter(mw.ValidateRouterID(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test/12345", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid router ID format")
}
