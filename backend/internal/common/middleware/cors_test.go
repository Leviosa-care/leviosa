package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnableCORS_DefaultOrigin(t *testing.T) {
	// allowedOrigin is "http://localhost:5173" by default; snapshot before any
	// test-local SetAllowedOrigin call changes it.
	want := allowedOrigin

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	called := false
	EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})(w, req)

	if !called {
		t.Fatal("next handler was not called")
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != want {
		t.Errorf("expected origin %s, got %s", want, origin)
	}
}

func TestEnableCORS_Headers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	EnableCORS(func(w http.ResponseWriter, r *http.Request) {})(w, req)

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods != "GET, POST, OPTIONS, PUT, DELETE" {
		t.Errorf("unexpected methods: %s", methods)
	}

	credentials := w.Header().Get("Access-Control-Allow-Credentials")
	if credentials != "true" {
		t.Errorf("expected credentials 'true', got %s", credentials)
	}
}

func TestEnableCORS_PreflightOptions(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	called := false
	EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})(w, req)

	if called {
		t.Fatal("next handler should not be called for OPTIONS preflight")
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin == "" {
		t.Error("expected Access-Control-Allow-Origin header to be set")
	}

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("expected Access-Control-Allow-Methods header to be set")
	}

	headers := w.Header().Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Error("expected Access-Control-Allow-Headers header to be set")
	}

	maxAge := w.Header().Get("Access-Control-Max-Age")
	if maxAge != "86400" {
		t.Errorf("expected Access-Control-Max-Age 86400, got %s", maxAge)
	}
}

func TestSetAllowedOrigin_SetsValue(t *testing.T) {
	// allowedOrigin starts at default; SetAllowedOrigin with sync.Once means
	// the *first* call wins. Since this is the first test that calls
	// SetAllowedOrigin in this process, it will succeed. We verify the
	// package-level var rather than retesting EnableCORS.
	SetAllowedOrigin("https://staging.leviosa.com")

	if allowedOrigin != "https://staging.leviosa.com" {
		t.Errorf("expected allowedOrigin https://staging.leviosa.com, got %s", allowedOrigin)
	}
}

func TestSetAllowedOrigin_IgnoresEmpty(t *testing.T) {
	before := allowedOrigin
	SetAllowedOrigin("")
	if allowedOrigin != before {
		t.Errorf("allowedOrigin should not change on empty, was %s now %s", before, allowedOrigin)
	}
}

func TestSetAllowedOrigin_SubsequentCallsIgnored(t *testing.T) {
	// sync.Once means only the first call takes effect.
	// (Already called in the test above, so this is a no-op.)
	SetAllowedOrigin("https://other.example.com")
	if allowedOrigin == "https://other.example.com" {
		t.Error("subsequent SetAllowedOrigin call should be ignored")
	}
}
