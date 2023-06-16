package integration

import (
	"encoding/json"
	"go-backend-sample/internal/handler"
	"testing"

	"github.com/google/uuid"
)

func TestUser(t *testing.T) {
	t.Run("create an user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test","email":"test@example.com"}`)
			assert(t, 200, rec.Code)

			res := handler.CreateUserResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(t, false, uuid.Nil == res.ID)

			userIDMap["user1"] = res.ID
		})

		t.Run("invalid json", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"name":"test","email":`)
			assert(t, 400, rec.Code)
		})

		t.Run("invalid request body", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "POST", "/api/v1/users", `{"email":"test2@example.com"}`)
			assert(t, 400, rec.Code)

			rec = doRequest(t, "POST", "/api/v1/users", `{"name":"test2"}`)
			assert(t, 400, rec.Code)

			rec = doRequest(t, "POST", "/api/v1/users", `{"name":"test2","email":"not_email"}`)
			assert(t, 400, rec.Code)
		})
	})

	t.Run("get users", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "GET", "/api/v1/users", "")
			assert(t, 200, rec.Code)

			res := handler.GetUsersResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(
				t,
				handler.GetUsersResponse{
					{
						ID:    userIDMap["user1"],
						Name:  "test",
						Email: "test@example.com",
					},
				},
				res,
			)
		})
	})

	t.Run("get an user", func(t *testing.T) {
		t.Run("success: user1", func(t *testing.T) {
			t.Parallel()
			rec := doRequest(t, "GET", "/api/v1/users/"+userIDMap["user1"].String(), "")
			assert(t, 200, rec.Code)

			res := handler.GetUserResponse{}
			assert(t, nil, json.Unmarshal(rec.Body.Bytes(), &res))
			assert(
				t,
				handler.GetUserResponse{
					ID:    userIDMap["user1"],
					Name:  "test",
					Email: "test@example.com",
				},
				res,
			)
		})
	})
}
