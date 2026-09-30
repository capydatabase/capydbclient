package capydbclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestDoer returns a Doer pointed at server with no retries configured, which
// is the shape both consumers use for non-GET calls.
func newTestDoer(server *httptest.Server) *Doer {
	return &Doer{
		APIKey:     "capy_live_test",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		UserAgent:  "capydb-test/1.0",
	}
}

func TestDoRequestShaping(t *testing.T) {
	t.Parallel()

	type capture struct {
		accept      string
		auth        string
		body        string
		contentType string
		method      string
		path        string
		userAgent   string
		extra       string
	}

	tests := []struct {
		name    string
		method  string
		path    string
		payload any
		extra   []Header
		want    capture
	}{
		{
			name:   "get carries auth and accept but no content type",
			method: http.MethodGet,
			path:   "/v1/projects",
			want: capture{
				accept:    "application/json",
				auth:      "Bearer capy_live_test",
				method:    http.MethodGet,
				path:      "/v1/projects",
				userAgent: "capydb-test/1.0",
			},
		},
		{
			name:    "post encodes payload and sets content type",
			method:  http.MethodPost,
			path:    "/v1/projects",
			payload: map[string]string{"name": "atlas"},
			want: capture{
				accept:      "application/json",
				auth:        "Bearer capy_live_test",
				body:        `{"name":"atlas"}`,
				contentType: "application/json",
				method:      http.MethodPost,
				path:        "/v1/projects",
				userAgent:   "capydb-test/1.0",
			},
		},
		{
			name:   "extra header is applied",
			method: http.MethodGet,
			path:   "/v1/cli/device/token",
			extra:  []Header{{Key: "X-Device-Token", Value: "secret"}},
			want: capture{
				accept:    "application/json",
				auth:      "Bearer capy_live_test",
				extra:     "secret",
				method:    http.MethodGet,
				path:      "/v1/cli/device/token",
				userAgent: "capydb-test/1.0",
			},
		},
		{
			// A blank key or value must not produce an empty header - the device
			// login poll passes a zero Header before the token exists.
			name:   "blank extra header is skipped",
			method: http.MethodGet,
			path:   "/v1/me",
			extra:  []Header{{Key: "X-Device-Token", Value: ""}, {Key: "", Value: "orphan"}},
			want: capture{
				accept:    "application/json",
				auth:      "Bearer capy_live_test",
				method:    http.MethodGet,
				path:      "/v1/me",
				userAgent: "capydb-test/1.0",
			},
		},
		{
			// The caller owns escaping; the Doer must pass the path through
			// byte-for-byte so an already-escaped segment is not double-encoded.
			name:   "pre-escaped path is passed through verbatim",
			method: http.MethodGet,
			path:   "/v1/projects/proj%2Fwith%20space/backups",
			want: capture{
				accept:    "application/json",
				auth:      "Bearer capy_live_test",
				method:    http.MethodGet,
				path:      "/v1/projects/proj%2Fwith%20space/backups",
				userAgent: "capydb-test/1.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got capture
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				raw, readErr := io.ReadAll(r.Body)
				if readErr != nil {
					t.Errorf("read request body: %v", readErr)
				}
				got = capture{
					accept:      r.Header.Get("Accept"),
					auth:        r.Header.Get("Authorization"),
					body:        string(raw),
					contentType: r.Header.Get("Content-Type"),
					method:      r.Method,
					// RequestURI preserves the raw, un-decoded path so a
					// double-escaping regression is visible.
					path:      r.RequestURI,
					userAgent: r.Header.Get("User-Agent"),
					extra:     r.Header.Get("X-Device-Token"),
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()

			doer := newTestDoer(server)
			if err := doer.Do(context.Background(), tt.method, tt.path, tt.payload, nil, tt.extra...); err != nil {
				t.Fatalf("Do() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("request shaping\n got: %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}

func TestDoOmitsAuthorizationWithoutAPIKey(t *testing.T) {
	t.Parallel()

	var present bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, present = r.Header["Authorization"]
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	doer := newTestDoer(server)
	doer.APIKey = ""
	if err := doer.Do(context.Background(), http.MethodGet, "/v1/health", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if present {
		t.Fatal("Authorization header must be absent when no API key is configured")
	}
}

func TestDoDecodesResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"proj_1","name":"atlas"}`))
	}))
	defer server.Close()

	var dest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := newTestDoer(server).Do(context.Background(), http.MethodGet, "/v1/projects/proj_1", nil, &dest); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if dest.ID != "proj_1" || dest.Name != "atlas" {
		t.Fatalf("decoded = %+v, want {proj_1 atlas}", dest)
	}
}

func TestDoDecodeEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		body        string
		wantDest    bool
		wantErrPart string
	}{
		{name: "nil dest discards body", body: `{"id":"proj_1"}`, wantDest: false},
		// 200 with an empty body is how the control plane answers some
		// mutations; it must not be reported as a decode failure.
		{name: "empty body with dest is not an error", body: "", wantDest: true},
		{name: "malformed json surfaces a decode error", body: `{"id":`, wantDest: true, wantErrPart: "decode response body"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			var dest map[string]any
			var target any
			if tt.wantDest {
				target = &dest
			}

			err := newTestDoer(server).Do(context.Background(), http.MethodGet, "/v1/x", nil, target)
			if tt.wantErrPart == "" {
				if err != nil {
					t.Fatalf("Do() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErrPart) {
				t.Fatalf("Do() error = %v, want error containing %q", err, tt.wantErrPart)
			}
		})
	}
}

func TestDoAPIError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       int
		body         string
		wantMessage  string
		wantText     string
		wantNotFound bool
	}{
		{
			name:        "error field becomes the message",
			status:      http.StatusConflict,
			body:        `{"error":"project has an active operation"}`,
			wantMessage: "project has an active operation",
			wantText:    "capydb api request failed with status 409: project has an active operation",
		},
		{
			name:     "missing error field falls back to status only",
			status:   http.StatusForbidden,
			body:     `{}`,
			wantText: "capydb api request failed with status 403",
		},
		{
			// An HTML error page from a proxy must not break error reporting.
			name:     "non json body still yields an APIError",
			status:   http.StatusBadGateway,
			body:     "<html>bad gateway</html>",
			wantText: "capydb api request failed with status 502",
		},
		{
			name:         "404 is detected by IsNotFound",
			status:       http.StatusNotFound,
			body:         `{"error":"project not found"}`,
			wantMessage:  "project not found",
			wantText:     "capydb api request failed with status 404: project not found",
			wantNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			// POST so the 5xx case is not retried by this assertion.
			err := newTestDoer(server).Do(context.Background(), http.MethodPost, "/v1/x", map[string]string{}, nil)

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("Do() error = %v, want *APIError", err)
			}
			if apiErr.StatusCode != tt.status {
				t.Fatalf("StatusCode = %d, want %d", apiErr.StatusCode, tt.status)
			}
			if apiErr.Message != tt.wantMessage {
				t.Fatalf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}
			if apiErr.Error() != tt.wantText {
				t.Fatalf("Error() = %q, want %q", apiErr.Error(), tt.wantText)
			}
			if IsNotFound(err) != tt.wantNotFound {
				t.Fatalf("IsNotFound() = %v, want %v", IsNotFound(err), tt.wantNotFound)
			}
		})
	}
}

func TestIsNotFoundRejectsNonAPIErrors(t *testing.T) {
	t.Parallel()

	if IsNotFound(nil) {
		t.Fatal("IsNotFound(nil) must be false")
	}
	if IsNotFound(errors.New("dial tcp: connection refused")) {
		t.Fatal("IsNotFound(plain error) must be false")
	}
	// Wrapped APIErrors must still be recognized - both consumers wrap with %w.
	wrapped := fmt.Errorf("get project: %w", &APIError{StatusCode: http.StatusNotFound})
	if !IsNotFound(wrapped) {
		t.Fatal("IsNotFound must unwrap a wrapped APIError")
	}
}

func TestDoRetryPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		method       string
		statuses     []int
		retryBackoff []time.Duration
		wantAttempts int32
		wantErr      bool
	}{
		{
			name:         "get retries a 5xx then succeeds",
			method:       http.MethodGet,
			statuses:     []int{http.StatusInternalServerError, http.StatusOK},
			retryBackoff: []time.Duration{0, time.Millisecond},
			wantAttempts: 2,
		},
		{
			// A 4xx is the server's final answer; retrying only burns budget.
			name:         "get does not retry a 4xx",
			method:       http.MethodGet,
			statuses:     []int{http.StatusBadRequest, http.StatusOK},
			retryBackoff: []time.Duration{0, time.Millisecond},
			wantAttempts: 1,
			wantErr:      true,
		},
		{
			// The core safety property: a non-GET is never replayed, so a
			// create/delete cannot be applied twice.
			name:         "non get is never retried even on 5xx",
			method:       http.MethodPost,
			statuses:     []int{http.StatusInternalServerError, http.StatusOK},
			retryBackoff: []time.Duration{0, time.Millisecond},
			wantAttempts: 1,
			wantErr:      true,
		},
		{
			name:         "retries are exhausted and the last error is returned",
			method:       http.MethodGet,
			statuses:     []int{http.StatusServiceUnavailable, http.StatusServiceUnavailable, http.StatusServiceUnavailable},
			retryBackoff: []time.Duration{0, time.Millisecond, time.Millisecond},
			wantAttempts: 3,
			wantErr:      true,
		},
		{
			name:         "no retry backoff means a single attempt",
			method:       http.MethodGet,
			statuses:     []int{http.StatusInternalServerError, http.StatusOK},
			wantAttempts: 1,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				n := int(attempts.Add(1)) - 1
				status := tt.statuses[len(tt.statuses)-1]
				if n < len(tt.statuses) {
					status = tt.statuses[n]
				}
				w.WriteHeader(status)
			}))
			defer server.Close()

			doer := newTestDoer(server)
			doer.RetryBackoff = tt.retryBackoff

			var payload any
			if tt.method != http.MethodGet {
				payload = map[string]string{"name": "atlas"}
			}
			err := doer.Do(context.Background(), tt.method, "/v1/x", payload, nil)

			if tt.wantErr && err == nil {
				t.Fatal("Do() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Do() error = %v, want nil", err)
			}
			if got := attempts.Load(); got != tt.wantAttempts {
				t.Fatalf("server attempts = %d, want %d", got, tt.wantAttempts)
			}
		})
	}
}

func TestDoRetriesNetworkErrors(t *testing.T) {
	t.Parallel()

	// A closed server produces a transport error rather than an HTTP status,
	// which is the other retryable class.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	client := server.Client()
	url := server.URL
	server.Close()

	doer := &Doer{
		BaseURL:      url,
		HTTPClient:   client,
		UserAgent:    "capydb-test/1.0",
		RetryBackoff: []time.Duration{0, time.Millisecond},
	}
	err := doer.Do(context.Background(), http.MethodGet, "/v1/x", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "perform request") {
		t.Fatalf("Do() error = %v, want a wrapped transport error", err)
	}
}

func TestDoHonorsContextCancellationDuringBackoff(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		// Cancel while the client is about to sleep before the retry.
		cancel()
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	doer := newTestDoer(server)
	doer.RetryBackoff = []time.Duration{0, 5 * time.Second}

	start := time.Now()
	err := doer.Do(ctx, http.MethodGet, "/v1/x", nil, nil)
	if err == nil {
		t.Fatal("Do() error = nil, want cancellation")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("Do() waited %s; cancellation must abort the backoff sleep", elapsed)
	}
	if got := attempts.Load(); got != 1 {
		t.Fatalf("server attempts = %d, want 1", got)
	}
}

func TestDoRejectsUnencodablePayload(t *testing.T) {
	t.Parallel()

	var called atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// A channel cannot be marshaled; the request must fail before any I/O.
	err := newTestDoer(server).Do(context.Background(), http.MethodPost, "/v1/x", make(chan int), nil)
	if err == nil || !strings.Contains(err.Error(), "encode request body") {
		t.Fatalf("Do() error = %v, want an encode error", err)
	}
	if called.Load() {
		t.Fatal("no request should be sent when the payload cannot be encoded")
	}
}

func TestNormalizeList(t *testing.T) {
	t.Parallel()

	// The bug this guards: Go nil slices marshal as JSON null, and the control
	// plane's list endpoints historically returned null for an empty list, which
	// crashed UI callers doing .map/.filter.
	var nilProjects []Project
	normalized := NormalizeList(nilProjects)
	if normalized == nil {
		t.Fatal("NormalizeList(nil) must return a non-nil slice")
	}
	if len(normalized) != 0 {
		t.Fatalf("NormalizeList(nil) length = %d, want 0", len(normalized))
	}
	encoded, err := json.Marshal(normalized)
	if err != nil {
		t.Fatalf("marshal normalized list: %v", err)
	}
	if string(encoded) != "[]" {
		t.Fatalf("normalized empty list marshals as %s, want []", encoded)
	}

	populated := NormalizeList([]Project{{ID: "proj_1"}})
	if len(populated) != 1 || populated[0].ID != "proj_1" {
		t.Fatalf("NormalizeList must preserve entries, got %+v", populated)
	}
}

func TestNormalizeListDecodedNull(t *testing.T) {
	t.Parallel()

	// End to end: a `null` list field decodes to nil and must normalize to [].
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"projects":null}`))
	}))
	defer server.Close()

	var body struct {
		Projects []Project `json:"projects"`
	}
	if err := newTestDoer(server).Do(context.Background(), http.MethodGet, "/v1/projects", nil, &body); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if body.Projects != nil {
		t.Fatal("a JSON null list is expected to decode to nil before normalization")
	}
	if got := NormalizeList(body.Projects); got == nil || len(got) != 0 {
		t.Fatalf("NormalizeList(decoded null) = %#v, want empty non-nil slice", got)
	}
}
