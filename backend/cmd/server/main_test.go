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
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("expected ok health response, got %s", rec.Body.String())
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
	if !strings.Contains(rec.Body.String(), "Tokyo Loop Night") {
		t.Fatalf("expected seeded event, got %s", rec.Body.String())
	}
}

func TestJoinValidation(t *testing.T) {
	server := NewServer()
	body := bytes.NewBufferString(`{"nickname":"Local Tester","email":"not-an-email"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/join", body)

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "email is invalid") {
		t.Fatalf("expected validation error, got %s", rec.Body.String())
	}
}

func TestJoinAccepted(t *testing.T) {
	server := NewServer()
	payload := JoinRequest{
		Nickname: "Local Tester",
		Email:    "local@example.com",
		Role:     "viewer",
		Message:  "Testing",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/join", bytes.NewReader(body))

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"accepted":true`) {
		t.Fatalf("expected accepted response, got %s", rec.Body.String())
	}
}
