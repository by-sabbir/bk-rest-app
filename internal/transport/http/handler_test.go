package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/by-sabbir/company-microservice-rest/internal/company"
	"github.com/by-sabbir/company-microservice-rest/internal/db"
	transportHttp "github.com/by-sabbir/company-microservice-rest/internal/transport/http"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestCompanyAPIPositive(t *testing.T) {
	t.Setenv("DB_HOST", "0.0.0.0")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USERNAME", "bktest")
	t.Setenv("DB_PASSWORD", "hello")
	t.Setenv("DB_NAME", "postgres")
	t.Setenv("SSL_MODE", "disable")
	t.Setenv("JWT_SECRET", "bk-go-dev")
	rand.Seed(time.Now().UnixNano())

	var id string
	createRequest := &company.Company{
		Name:           RandStringBytes(10),
		Description:    "Lorem Ipsum Dolor Sit",
		TotalEmployees: 10,
		IsRegistered:   true,
		Type:           company.CompanyType[3],
	}

	t.Run("test healthcheck", func(t *testing.T) {
		uri := "/healthcheck"
		req := httptest.NewRequest("GET", uri, nil)
		resp := execReq(req)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	})
	t.Run("test create company api", func(t *testing.T) {
		uri := "/api/v1/private/company/create"
		payload, err := json.Marshal(createRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", uri, bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")

		resp := execReq(req)
		assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

		var got *company.Company
		err = json.Unmarshal(resp.Body.Bytes(), &got)
		assert.NoError(t, err)
		id = got.ID

		createRequest.ID = id
		assert.Equal(t, createRequest, got)
	})
	t.Run("test get company by id api", func(t *testing.T) {
		uri := fmt.Sprintf("/api/v1/public/company/%s", id)
		req := httptest.NewRequest("GET", uri, nil)
		resp := execReq(req)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		var got *company.Company
		err := json.Unmarshal(resp.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, id, got.ID)
	})
	t.Run("test partial update by id api", func(t *testing.T) {
		uri := fmt.Sprintf("/api/v1/private/company/patch/%s", id)
		updateRequest := company.Company{
			Type: company.CompanyType[3],
		}
		payload, err := json.Marshal(updateRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest("PATCH", uri, bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")

		resp := execReq(req)
		assert.Equal(t, http.StatusPartialContent, resp.Result().StatusCode)

		var got *company.Company
		err = json.Unmarshal(resp.Body.Bytes(), &got)
		assert.NoError(t, err)

		assert.Equal(t, updateRequest.Type, got.Type)
	})
	t.Run("test delete company by id api", func(t *testing.T) {
		uri := fmt.Sprintf("/api/v1/private/company/delete/%s", id)
		req := httptest.NewRequest("DELETE", uri, nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")
		resp := execReq(req)
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	})
}
func TestCompanyAPINegative(t *testing.T) {
	t.Setenv("DB_HOST", "0.0.0.0")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USERNAME", "bktest")
	t.Setenv("DB_PASSWORD", "hello")
	t.Setenv("DB_NAME", "postgres")
	t.Setenv("SSL_MODE", "disable")
	t.Setenv("JWT_SECRET", "bk-go-dev")
	rand.Seed(time.Now().UnixNano())
	uri := "/api/v1/private/company/create"

	t.Run("test create company api without name", func(t *testing.T) {

		createRequest := &company.Company{
			Description:    "Lorem Ipsum Dolor Sit",
			TotalEmployees: 10,
			IsRegistered:   true,
			Type:           company.CompanyType[3],
		}
		payload, err := json.Marshal(createRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", uri, bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")

		resp := execReq(req)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})
	t.Run("test create company api without type", func(t *testing.T) {

		createRequest := &company.Company{
			Name:           RandStringBytes(10),
			Description:    "Lorem Ipsum Dolor Sit",
			TotalEmployees: 10,
			IsRegistered:   true,
		}
		payload, err := json.Marshal(createRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", uri, bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")

		resp := execReq(req)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})
	t.Run("test create company api without total employee", func(t *testing.T) {

		createRequest := &company.Company{
			Name:         RandStringBytes(10),
			Description:  "Lorem Ipsum Dolor Sit",
			Type:         company.CompanyType[1],
			IsRegistered: true,
		}
		payload, err := json.Marshal(createRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", uri, bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")

		resp := execReq(req)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})
	t.Run("test get company by random id api", func(t *testing.T) {
		randomId := uuid.NewString()
		uri := fmt.Sprintf("/api/v1/public/company/%s", randomId)
		req := httptest.NewRequest("GET", uri, nil)
		resp := execReq(req)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})
	t.Run("test delete company by random id api", func(t *testing.T) {
		randomId := uuid.NewString()
		uri := fmt.Sprintf("/api/v1/private/company/delete/%s", randomId)
		req := httptest.NewRequest("DELETE", uri, nil)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")
		resp := execReq(req)
		assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
	})
}

func TestCompanyAPIPayload(t *testing.T) {
	t.Setenv("DB_HOST", "0.0.0.0")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USERNAME", "bktest")
	t.Setenv("DB_PASSWORD", "hello")
	t.Setenv("DB_NAME", "postgres")
	t.Setenv("SSL_MODE", "disable")
	t.Setenv("JWT_SECRET", "bk-go-dev")

	t.Run("test post payload", func(t *testing.T) {
		uri := "/api/v1/private/company/create"
		wrong_json_payload := strings.NewReader(`{
			"name": "Evil Corp",
			"description": "lorem ipsum dolor",
			"total_employees": 219,
			"type": "NonProfit"
			"is_registered": true
		  }`)
		req := httptest.NewRequest("POST", uri, wrong_json_payload)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")
		resp := execReq(req)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("test partial update payload", func(t *testing.T) {
		uri := "/api/v1/private/company/patch/41866014-e553-457c-89d8-ffa9b4114a3e"
		wrong_json_payload := strings.NewReader(`{
			"name": "Evil Corp",
			"is_registered": true`)
		req := httptest.NewRequest("PATCH", uri, wrong_json_payload)
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.aPo_C9u29lF0od_vU1V4Oox-OkVZGZMjKA2m_Wpn-D4")
		req.Header.Add("Content-Type", "application/json")
		resp := execReq(req)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

}

func TestCompanyAPIWithoutToken(t *testing.T) {
	t.Setenv("DB_HOST", "0.0.0.0")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USERNAME", "bktest")
	t.Setenv("DB_PASSWORD", "hello")
	t.Setenv("DB_NAME", "postgres")
	t.Setenv("SSL_MODE", "disable")
	t.Setenv("JWT_SECRET", "bk-go-dev")
	rand.Seed(time.Now().UnixNano())

	createRequest := &company.Company{
		Name:           RandStringBytes(10),
		Description:    "Lorem Ipsum Dolor Sit",
		TotalEmployees: 19,
		IsRegistered:   true,
		Type:           company.CompanyType[2],
	}
	t.Run("test create company api", func(t *testing.T) {
		uri := "/api/v1/private/company/create"
		payload, err := json.Marshal(createRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", uri, bytes.NewBuffer(payload))

		resp := execReq(req)
		assert.Equal(t, http.StatusUnauthorized, resp.Result().StatusCode)

	})
}

func execReq(req *http.Request) *httptest.ResponseRecorder {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error initiating test database: %+v\n", err)
	}
	svc := transportHttp.CompanyRestService(db)
	h := transportHttp.NewHandler(svc)
	rr := httptest.NewRecorder()

	h.Router.ServeHTTP(rr, req)
	return rr
}
