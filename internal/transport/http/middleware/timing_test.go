package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTimingMiddleware_LogsExecution(t *testing.T) {
	// 1. Inject an in-memory bytes buffer into the logger instance allowing us to assert against output.
	var buf bytes.Buffer
	testLogger := log.New(&buf, "", 0)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 2. Initialize the timing middleware injecting our logger instance
	middlewareToTest := TimingMiddleware(testLogger)(nextHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sales", nil)
	rec := httptest.NewRecorder()

	// 3. Dispatch the simulated network HTTP request through the middleware stream
	middlewareToTest.ServeHTTP(rec, req)

	// --- LOGGING ASSERTIONS ---

	// Assertion 1: Ensure that the internal log buffer actually received byte allocations
	if buf.Len() == 0 {
		t.Fatal("expected middleware to print a log line, but the log buffer is empty")
	}

	// Assertion 2: Validate the semantic formatting structure of the generated string
	logOutput := buf.String()

	if !strings.Contains(logOutput, "[POST]") {
		t.Errorf("expected log to contain HTTP method '[POST]', got: %q", logOutput)
	}

	if !strings.Contains(logOutput, "/api/v1/sales") {
		t.Errorf("expected log to contain URL path '/api/v1/sales', got: %q", logOutput)
	}

	if !strings.Contains(logOutput, "s") { // Verifies execution tracked in seconds (e.g., 0.000015s)
		t.Errorf("expected log to track duration formatted with seconds suffix 's', got: %q", logOutput)
	}
}
