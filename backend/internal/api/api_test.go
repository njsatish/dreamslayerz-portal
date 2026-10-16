package api

import (
	appconfig "github.com/njsatish/dreamslayerz-portal/backend/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	res := httptest.NewRecorder()
	cfg := appconfig.Config{}
	New(cfg).ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d", res.Code)
	}
	if !strings.Contains(res.Body.String(), "booksy-portal-api") {
		t.Fatalf("unexpected body: %s", res.Body.String())
	}
}
