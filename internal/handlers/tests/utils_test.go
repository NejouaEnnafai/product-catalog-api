package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// executeRequest creates a new ResponseRecorder, executes the request against the handler,
// and returns the response recorder
func executeRequest(req *http.Request, handler http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// checkResponseCode checks if the status code is as expected
func checkResponseCode(t testing.TB, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// parseResponse parses the JSON response into the given interface
func parseResponse(rr *httptest.ResponseRecorder, v interface{}) error {
	return json.Unmarshal(rr.Body.Bytes(), v)
}
