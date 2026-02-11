package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(handler gin.HandlerFunc) *httptest.ResponseRecorder {
	r := gin.New()
	r.GET("/", handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := performRequest(func(c *gin.Context) {
		Success(c, "Operation successful", map[string]string{"key": "value"})
	})

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"status":"success","message":"Operation successful","data":{"key":"value"}}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := performRequest(func(c *gin.Context) {
		Error(c, http.StatusBadRequest, "Bad request error")
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedBody := `{"status":"error","message":"Bad request error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
