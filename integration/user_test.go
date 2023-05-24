package integration

import (
	"encoding/json"
	"go-backend-sample/internal/handler"
	"go-backend-sample/internal/pkg/config"
	"go-backend-sample/internal/repository"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest/v3"
)

var (
	db        *sqlx.DB
	e         *echo.Echo
	r         *repository.Repository
	h         *handler.Handler
	userIDMap = make(map[string]uuid.UUID)
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("connect to docker: ", err)
	}

	if err := pool.Client.Ping(); err != nil {
		log.Fatal("ping docker: ", err)
	}

	mysqlConfig := config.MySQL()

	resource, err := pool.Run("mysql", "latest", []string{
		"MYSQL_ROOT_PASSWORD=" + mysqlConfig.Passwd,
		"MYSQL_DATABASE=" + mysqlConfig.DBName,
	})
	if err != nil {
		log.Fatal("run docker: ", err)
	}

	mysqlConfig.Addr = "localhost:" + resource.GetPort("3306/tcp")

	if err := pool.Retry(func() error {
		_db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
		if err != nil {
			return err
		}
		db = _db
		return _db.Ping()
	}); err != nil {
		log.Fatal("connect to database container: ", err)
	}

	// setup dependencies
	r = repository.New(db)
	if err := r.SetupTables(); err != nil {
		log.Fatal("setup tables: ", err)
	}
	h = handler.New(r)
	e = echo.New()
	h.SetupRoutes(e.Group("/api"))

	log.Println("start integration test")
	m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatal("purge docker: ", err)
	}
}

func doRequest(method, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestUser(t *testing.T) {
	t.Run("create an user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest("POST", "/api/users", `{"name":"test","email":"test@example.com"}`)
			if rec.Code != 200 {
				t.Errorf("expected: %v, actual: %v", 200, rec.Code)
			}

			res := handler.CreateUserResponse{}
			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Fatal("unmarshal response body: ", err)
			}

			if res.ID == uuid.Nil {
				t.Errorf("expected: %v, actual: %v", uuid.Nil, res.ID)
			}

			userIDMap["user1"] = res.ID
		})

		t.Run("invalid json", func(t *testing.T) {
			t.Parallel()
			rec := doRequest("POST", "/api/users", `{"name":"test","email":`)
			if rec.Code != 400 {
				t.Errorf("expected: %v, actual: %v", 400, rec.Code)
			}
		})

		t.Run("invalid request body", func(t *testing.T) {
			t.Parallel()
			rec := doRequest("POST", "/api/users", `{"email":"test2@example.com"}`)
			if rec.Code != 400 {
				t.Errorf("expected: %v, actual: %v", 400, rec.Code)
			}

			rec = doRequest("POST", "/api/users", `{"name":"test2"}`)
			if rec.Code != 400 {
				t.Errorf("expected: %v, actual: %v", 400, rec.Code)
			}

			rec = doRequest("POST", "/api/users", `{"name":"test2","email":"not_email"}`)
			if rec.Code != 400 {
				t.Errorf("expected: %v, actual: %v", 400, rec.Code)
			}
		})
	})

	t.Run("get users", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			rec := doRequest("GET", "/api/users", "")
			if !cmp.Equal(200, rec.Code) {
				t.Errorf("expected: %v, actual: %v", 200, rec.Code)
			}

			res := handler.GetUsersResponse{}
			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Fatal("unmarshal response body: ", err)
			}

			diff := cmp.Diff(res, handler.GetUsersResponse{
				{
					ID:    userIDMap["user1"],
					Name:  "test",
					Email: "test@example.com",
				},
			})
			if len(diff) > 0 {
				t.Errorf("diff: %v", diff)
			}
		})
	})

	t.Run("get an user", func(t *testing.T) {
		t.Run("success: user1", func(t *testing.T) {
			t.Parallel()
			rec := doRequest("GET", "/api/users/"+userIDMap["user1"].String(), "")
			if rec.Code != 200 {
				t.Errorf("expected: %v, actual: %v", 200, rec.Code)
			}

			res := handler.GetUserResponse{}
			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Fatal("unmarshal response body: ", err)
			}

			diff := cmp.Diff(res, handler.GetUserResponse{
				ID:    userIDMap["user1"],
				Name:  "test",
				Email: "test@example.com",
			})
			if len(diff) > 0 {
				t.Errorf("diff: %v", diff)
			}
		})
	})
}
