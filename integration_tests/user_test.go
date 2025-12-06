// NOTE: go test -updateを実行することで、スナップショットを更新することができる

package integrationtests

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

		t.Run("invalid: name is blank", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"email":"test2@example.com"}`)

			expectedStatus := `400 Bad Request`
			expectedBody := `{"error_message":"operation CreateUser: decode request: decode application/json: invalid: name (field required)"}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})

		t.Run("invalid: email is blank", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test2"}`)

			expectedStatus := `400 Bad Request`
			expectedBody := `{"error_message":"operation CreateUser: decode request: decode application/json: invalid: email (field required)"}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})

		t.Run("invalid: email is invalid", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test2","email":"not_email"}`)

			expectedStatus := `400 Bad Request`
			expectedBody := `{"error_message":"operation CreateUser: decode request: validate: invalid: email (string: no @)"}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})
	})

	t.Run("get users", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "GET", "/api/v1/users", "")

			expectedStatus := `200 OK`
			expectedBody := `[{"id":"[UUID]","name":"test","email":"test@example.com","iconUrl":"https://via.placeholder.com/150/92c952"}]`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})
	})

	t.Run("get an user", func(t *testing.T) {
		t.Run("success: user1", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "GET", "/api/v1/users/"+userIDMap["user1"].String(), "")

			expectedStatus := `200 OK`
			expectedBody := `{"id":"[UUID]","name":"test","email":"test@example.com","iconUrl":"https://via.placeholder.com/150/92c952"}`
			assert.Equal(t, rec.Result().Status, expectedStatus)
			assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
		})
	})
}
