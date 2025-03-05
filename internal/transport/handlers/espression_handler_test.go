package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/storage"
)

func TestExpressionHandler_ServeHTTP(t *testing.T) {
	s := storage.NewStorage()
	exprService := expression.NewExpressionService(s, timeConfig)
	handler := NewExpressionHandler(exprService)

	t.Run("Valid ID", func(t *testing.T) {
		exprService.ProcessExpression("2+2")

		req, err := http.NewRequest(http.MethodGet, "/api/v1/expression/1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Handle("/api/v1/expression/{id}", handler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NoError(t, err)
	})

	t.Run("Missing ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/expression/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Handle("/api/v1/expression/{id}", handler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/expression/abc", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Handle("/api/v1/expression/{id}", handler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "id must be a number", response["error"])
	})
}
