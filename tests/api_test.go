package tests

import (
	"dds-backend/routes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO make tests
func TestPingRoute(t *testing.T) {
	router, _, _, err := routes.MakeServer()

	assert.Equal(t, nil, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 418, w.Code)
	//var response map[string]string
	//json.Unmarshal([]byte(w.Body.String()), &response)

}
