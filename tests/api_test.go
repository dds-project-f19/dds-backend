package tests

import (
	"github.com/gin-gonic/gin"
	//"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := gin.Default()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	router.ServeHTTP(w, req)

	//assert.Equal(t, 200, w.Code)
	//assert.Equal(t, "PONG", w.Body.String())

}
