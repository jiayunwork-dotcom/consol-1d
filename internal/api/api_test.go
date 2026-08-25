package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConsolidateEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"scenario":{"cv":1e-7,"thickness":10,"drainage":"double","initial_pressure":{"type":"uniform","u0":100},"time":5e7},"nodes":21}`))
	req := httptest.NewRequest(http.MethodPost, "/api/consolidate", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out struct {
		U float64 `json:"u"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.U < 0.45 || out.U > 0.55 {
		t.Fatalf("U=%v", out.U)
	}
	if s.Book().Len() != 1 {
		t.Fatalf("book len=%d", s.Book().Len())
	}
}

func TestCurveEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"scenario":{"cv":1e-7,"thickness":10,"drainage":"double","initial_pressure":{"type":"uniform","u0":100},"time":1},"times":[1e6,5e7],"nodes":11}`))
	req := httptest.NewRequest(http.MethodPost, "/api/curve", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("curve len=%d", len(out))
	}
}

func TestSettleEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"scenario":{"cv":1e-7,"thickness":10,"drainage":"double","initial_pressure":{"type":"uniform","u0":100},"time":1},"target":0.9}`))
	req := httptest.NewRequest(http.MethodPost, "/api/settle", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "time_s") {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}
