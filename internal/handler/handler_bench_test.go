package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkHealth(b *testing.B) {
	handler := New(nil, nil).Routes()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()

	for b.Loop() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
	}
}
