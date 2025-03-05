package models

import (
	"time"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/calculation"
)

type Expression struct {
	ID         int               `json:"id"`
	Status     string            `json:"status"`
	Result     float64           `json:"result"`
	BinaryTree *calculation.Tree `json:"-"`
}

type Task struct {
	ID            int           `json:"id"`
	ExpressionID  int           `json:"-"`
	Status        string        `json:"-"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}
