package main

import (
	"context"
	"net/http"
	"testing"
)

func TestFailureClassification(t *testing.T) {
	tests := []struct {
		status int
		code   string
		want   string
	}{
		{http.StatusTooManyRequests, "ROUTING_BUSY", "routing_busy"},
		{http.StatusGatewayTimeout, "ROUTING_TIMEOUT", "timeout"},
		{http.StatusConflict, "NO_ELIGIBLE_DEVICE", "conflict"},
		{http.StatusUnprocessableEntity, "NO_ROUTE", "no_route"},
		{http.StatusBadRequest, "INVALID_REQUEST", "client_error"},
		{http.StatusInternalServerError, "ORCHESTRATION_ERROR", "server_error"},
	}
	for _, test := range tests {
		if got := classifyHTTP(test.status, test.code); got != test.want {
			t.Fatalf("HTTP %d %s classified as %s, want %s", test.status, test.code, got, test.want)
		}
	}
	if got := classifyTransport(context.Canceled); got != "cancelled" {
		t.Fatalf("cancel classified as %s", got)
	}
	if got := classifyTransport(context.DeadlineExceeded); got != "timeout" {
		t.Fatalf("deadline classified as %s", got)
	}
}

func TestUnexpectedTotalsExcludeBoundedFailures(t *testing.T) {
	values := map[string]int64{"route.routing_busy": 5, "task.conflict": 2, "route.timeout": 1, "nearby.server_error": 0, "command.transport_error": 0, "task.unexpected": 0}
	totals := errorTotals(values)
	if totals["expected"] != 8 || totals["unexpected"] != 0 {
		t.Fatalf("unexpected totals: %#v", totals)
	}
	values["nearby.server_error"] = 1
	if errorTotals(values)["unexpected"] != 1 {
		t.Fatal("server error did not fail the unexpected-error gate")
	}
}
