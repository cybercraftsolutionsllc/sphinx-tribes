package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGithubClientWithTokenDoesNotSendEmptyAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("expected no Authorization header for empty token, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"title":"Public issue","state":"open","body":"body"}`))
	}))
	defer server.Close()

	client := githubClientWithToken("")
	baseURL, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatalf("failed to parse test server url: %v", err)
	}
	client.BaseURL = baseURL

	if _, _, err := client.Issues.Get(context.Background(), "owner", "repo", 1); err != nil {
		t.Fatalf("expected public issue request without a token to succeed: %v", err)
	}
}

func TestGithubClientWithTokenSendsAuthorizationWhenConfigured(t *testing.T) {
	const token = "test-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			t.Fatalf("expected bearer Authorization header, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"title":"Private issue","state":"open","body":"body"}`))
	}))
	defer server.Close()

	client := githubClientWithToken(token)
	baseURL, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatalf("failed to parse test server url: %v", err)
	}
	client.BaseURL = baseURL

	if _, _, err := client.Issues.Get(context.Background(), "owner", "repo", 1); err != nil {
		t.Fatalf("expected issue request with a token to succeed: %v", err)
	}
}
