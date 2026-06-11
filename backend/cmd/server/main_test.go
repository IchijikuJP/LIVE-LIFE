package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"livelife/backend/internal/application"
	"livelife/backend/internal/domain"
	sqlitestore "livelife/backend/internal/infrastructure/sqlite"
	"livelife/backend/internal/interfaces/httpapi"
)

type testHarness struct {
	server http.Handler
	store  *sqlitestore.Store
}

func newTestHarness(t *testing.T) testHarness {
	t.Helper()

	store, err := sqlitestore.NewStore(":memory:")
	if err != nil {
		t.Fatalf("create sqlite store: %v", err)
	}
	service := application.NewService(store)
	server := httpapi.NewServer(service, filepath.Join("..", "..", "static"))

	return testHarness{server: server, store: store}
}

func TestHealth(t *testing.T) {
	harness := newTestHarness(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)

	harness.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("expected ok health response, got %s", body)
	}
	if !strings.Contains(body, `"brand":"LIVE LIFE"`) {
		t.Fatalf("expected LIVE LIFE brand in health response, got %s", body)
	}
}

func TestEvents(t *testing.T) {
	harness := newTestHarness(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil)

	harness.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "LIVE LIFE presents") {
		t.Fatalf("expected LIVE LIFE owned event, got %s", body)
	}
	if !strings.Contains(body, "ownedEvents") {
		t.Fatalf("expected ownedEvents group, got %s", body)
	}
	if !strings.Contains(body, "紅髪少年殺人事件") {
		t.Fatalf("expected clean Japanese/Chinese seed text, got %s", body)
	}
}

func TestCDItems(t *testing.T) {
	harness := newTestHarness(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/cd-items", nil)

	harness.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{`"brand":"LIVE LIFE"`, `"cd"`, `"vinyl"`, `"purchaseUrl"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected %s in CD response, got %s", want, body)
		}
	}
}

func TestShopItemsNoLongerTopLevel(t *testing.T) {
	harness := newTestHarness(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/shop-items", nil)

	harness.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d for removed top-level shop API, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestConnectValidation(t *testing.T) {
	harness := newTestHarness(t)
	body := bytes.NewBufferString(`{"nickname":"Local Tester","email":"not-an-email"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/connect", body)

	harness.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "email is invalid") {
		t.Fatalf("expected validation error, got %s", rec.Body.String())
	}
}

func TestConnectAccepted(t *testing.T) {
	harness := newTestHarness(t)
	payload := domain.ConnectRequest{
		Nickname: "Local Tester",
		Email:    "local@example.com",
		Topic:    "cd-select",
		Message:  "Testing",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/connect", bytes.NewReader(body))

	harness.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}
	response := rec.Body.String()
	if !strings.Contains(response, `"accepted":true`) {
		t.Fatalf("expected accepted response, got %s", response)
	}
	if !strings.Contains(response, `"brand":"LIVE LIFE"`) {
		t.Fatalf("expected LIVE LIFE brand in accepted response, got %s", response)
	}

	count, err := harness.store.CountConnectMessages(context.Background())
	if err != nil {
		t.Fatalf("expected connect message query to succeed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one persisted connect message, got %d", count)
	}
}
