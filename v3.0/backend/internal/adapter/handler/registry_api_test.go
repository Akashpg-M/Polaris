package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDIsStableWithinRequest(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	first := RequestID(c)
	if first == "" || RequestID(c) != first {
		t.Fatal("generated request ID changed within one request")
	}
}

func TestRequestIDPreservesCallerValue(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Request-ID", "caller-request")
	if got := RequestID(c); got != "caller-request" || RequestID(c) != got {
		t.Fatalf("request ID = %q, want caller-request", got)
	}
}
