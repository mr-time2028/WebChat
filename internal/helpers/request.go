package helpers

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func FakePostRequest(body string, url string, handler http.Handler) *httptest.ResponseRecorder {
	reader := strings.NewReader(body)
	req, _ := http.NewRequest("POST", url, reader)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}
