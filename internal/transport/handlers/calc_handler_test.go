package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/config"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/storage"
)

var timeConfig = config.TimeConfig{
	TimeAdd: time.Second,
	TimeSub: time.Second,
	TimeMul: time.Second,
	TimeDiv: time.Second,
}

func TestCalcHandler_ServeHTTP(t *testing.T) {
	s := storage.NewStorage()
	exprService := expression.NewExpressionService(s, timeConfig)
	handler := NewCalcHandler(exprService)

	requestBody := map[string]string{
		"expression": "3 + 4",
	}
	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
