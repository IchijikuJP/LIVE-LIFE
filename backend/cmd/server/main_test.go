package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	server := NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)

	server.ServeHTTP(rec, req)

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
	server := NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/events", nil)

	server.ServeHTTP(rec, req)

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
}

func TestCDItems(t *testing.T) {
	server := NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/cd-items", nil)

	server.ServeHTTP(rec, req)

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
	server := NewServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/shop-items", nil)

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d for removed top-level shop API, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestConnectValidation(t *testing.T) {
	server := NewServer()
	body := bytes.NewBufferString(`{"nickname":"Local Tester","email":"not-an-email"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/connect", body)

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "email is invalid") {
		t.Fatalf("expected validation error, got %s", rec.Body.String())
	}
}

func TestConnectAccepted(t *testing.T) {
	server := NewServer()
	payload := ConnectRequest{
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

	server.ServeHTTP(rec, req)

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

	var count int64
	if err := server.db.Model(&ConnectMessageModel{}).Count(&count).Error; err != nil {
		t.Fatalf("expected connect message query to succeed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one persisted connect message, got %d", count)
	}
}
