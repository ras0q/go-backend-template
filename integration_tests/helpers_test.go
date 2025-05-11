package integration_tests

import (
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"gotest.tools/v3/assert"
)

var (
	uuidRegexp = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	timeRegexp = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{3}Z`)
)

func escapeSnapshot(t *testing.T, s string) string {
	t.Helper()

	s = strings.Trim(s, "\n")
	s = uuidRegexp.ReplaceAllString(s, "[UUID]")
	s = timeRegexp.ReplaceAllString(s, "[TIME]")

	return s
}

func doRequest(t *testing.T, method, path string, bodystr string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(bodystr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}

func unmarshalResponse(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	v := map[string]any{}
	assert.NilError(t, json.Unmarshal(rec.Body.Bytes(), &v))

	return v
}
