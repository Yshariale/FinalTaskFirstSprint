package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/storage"
)

func TestTaskHandler_GiveTask(t *testing.T) {
	s := storage.NewStorage()
	exprService := expression.NewExpressionService(s, timeConfig)
	handler := NewTaskHandler(exprService)

	req, err := http.NewRequest(http.MethodGet, "/internal/task", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestTaskHandler_ReceiveTask(t *testing.T) {
	s := storage.NewStorage()
	exprService := expression.NewExpressionService(s, timeConfig)
	exprService.ProcessExpression("2+2")
	handler := NewTaskHandler(exprService)

	task := map[string]interface{}{
		"id":     1,
		"result": 4,
	}
	body, err := json.Marshal(task)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/internal/task", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
