package web_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/web"
	"github.com/stretchr/testify/assert"
)

var uuidInDb, uuidNotInDb uuid.UUID

func setUp() *gin.Engine {
	uuidInDb = uuid.New()
	uuidNotInDb = uuid.New()

	database.Instance = database.NewInMemoryDb()
	database.Instance.Create(database.New(uuidInDb))

	appCfg := config.AppCfg{
		Debug: true,
	}
	return web.NewRouter(appCfg)

}

func TestExpenseRoute(t *testing.T) {
	router := setUp()
	type testCase struct {
		name string
		uuid string
		code int
	}

	tcs := []testCase{
		{
			name: "No UUID",
			uuid: "",
			code: http.StatusNotFound,
		},
		{
			name: "Invalid UUID",
			uuid: "not-a-valid-uuid-not-even-one-bit",
			code: http.StatusBadRequest,
		},
		{
			name: "Non-existant UUID",
			uuid: uuidNotInDb.String(),
			code: http.StatusNotFound,
		},
		{
			name: "Valid and existing UUID",
			uuid: uuidInDb.String(),
			code: http.StatusOK,
		},
	}

	for _, tc := range tcs {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/expense/%s", tc.uuid), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, tc.code, w.Code, tc.name)
		// Could also test (un)marshal to test if data is returned properly without missing fields, corruption...
	}

}
