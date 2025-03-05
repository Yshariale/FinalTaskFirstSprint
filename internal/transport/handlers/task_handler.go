package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
)

type TaskHandler struct {
	expressionService *expression.ExpressionService
}

func NewTaskHandler(expressionService *expression.ExpressionService) *TaskHandler {
	return &TaskHandler{
		expressionService: expressionService,
	}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.giveTask(w)
	} else if r.Method == http.MethodPost {
		h.receiveTask(w, r)
	}
}

func (h *TaskHandler) giveTask(w http.ResponseWriter) {
	// Логика спрятана сюда
	task := h.expressionService.GetPendingTask()
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "no tasks available"})
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) receiveTask(w http.ResponseWriter, r *http.Request) {
	var agent_request struct {
		TaskID int     `json:"id"`
		Result float64 `json:"result"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&agent_request); err != nil {
		slog.Error("Agent sent an invalid body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid body"})
		return
	}
	// Логика спрятана сюда
	h.expressionService.ProcessIncomingTask(agent_request.TaskID, agent_request.Result)
}
