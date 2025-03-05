package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/storage"
)

func TestExpressionListHandler_ServeHTTP(t *testing.T) {
	s := storage.NewStorage()
	exprService := expression.NewExpressionService(s, timeConfig)
	handler := NewExpressionListHandler(exprService)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}
