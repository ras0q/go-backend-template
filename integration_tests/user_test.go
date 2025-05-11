// NOTE: go test -updateを実行することで、スナップショットを更新することができる

package integration_tests

import (
	"testing"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestUser(t *testing.T) {
	var userIDMap = make(map[string]uuid.UUID)

	t.Run("create an user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test","email":"test@example.com"}`)

			expectedStatus := `200 OK`
			expectedBody := `{"id":"[UUID]"}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)

			res := unmarshalResponse(t, rec)
			userIDMap["user1"] = uuid.MustParse(res["id"].(string))
		})

		t.Run("invalid: email is blank", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"email":"test2@example.com"}`)

			expectedStatus := `400 Bad Request`
			expectedBody := `{"message":"invalid request body: name: cannot be blank."}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})

		t.Run("invalid: name is blank", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test2"}`)

			expectedStatus := `400 Bad Request`
			expectedBody := `{"message":"invalid request body: email: cannot be blank."}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})

		t.Run("invalid: email is invalid", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test2","email":"not_email"}`)

			expectedStatus := `400 Bad Request`
			expectedBody := `{"message":"invalid request body: email: must be a valid email address."}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})
	})

	t.Run("get users", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "GET", "/api/v1/users", "")

			expectedStatus := `200 OK`
			expectedBody := `[{"id":"[UUID]","name":"test","email":"test@example.com"}]`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})
	})

	t.Run("get an user", func(t *testing.T) {
		t.Run("success: user1", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "GET", "/api/v1/users/"+userIDMap["user1"].String(), "")

			expectedStatus := `200 OK`
			expectedBody := `{"id":"[UUID]","name":"test","email":"test@example.com"}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})
	})
}
