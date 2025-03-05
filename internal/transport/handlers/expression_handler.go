package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
	"github.com/gorilla/mux"
)

type ExpressionHandler struct {
	expressionService *expression.ExpressionService
}

func NewExpressionHandler(expressionService *expression.ExpressionService) *ExpressionHandler {
	return &ExpressionHandler{
		expressionService: expressionService,
	}
}

func (h *ExpressionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expression_id_str := vars["id"]

	if expression_id_str == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id is required"})
		return
	}

	expression_id, err := strconv.Atoi(expression_id_str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id must be a number"})
		return
	}

	// Логика спрятана сюда
	expression := h.expressionService.GetExpressionByID(expression_id)

	if expression == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "expression not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expression)
}
