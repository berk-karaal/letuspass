package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/router"
)

func TestHandleMetricsStatus(t *testing.T) {
	type Response struct {
		Status string `json:"status"`
	}

	apiConfig := config.NewRestapiConfigFromEnv()
	r := router.SetupRouter(apiConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/metrics/status", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("expected json response; got %s", w.Body.String())
	}

	if resp.Status != "OK" {
		t.Errorf("expected status OK; got %s", resp.Status)
	}
}
