package storage

//
// В этом модуле низкоуровневая логика взаимодействия
// со списком задач и выражений
//
// Все методы понятны и без моих комментариев
//

import (
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/models"
)

type Storage struct {
	mu             sync.Mutex
	taskList       []models.Task
	expressionList []models.Expression
	expressionID   int64
	taskID         int64
}

func NewStorage() *Storage {
	return &Storage{
		taskList:       make([]models.Task, 0),
		expressionList: make([]models.Expression, 0),
	}
}

func (s *Storage) AddExpression(expression models.Expression) int {
	slog.Info("Storage.AddExpression: Получено выражение", "expression", expression)
	s.mu.Lock()
	defer s.mu.Unlock()

	expression.ID = int(atomic.AddInt64(&s.expressionID, 1))
	slog.Info("Storage.AddExpression: Выражению назначен ID", "expression", expression)
	s.expressionList = append(s.expressionList, expression)
	slog.Info("Storage.AddExpression: Выражение добавлено в список", "expression", expression, "s.expressionList", s.expressionList)
	return expression.ID
}

func (s *Storage) GetExpressions() []models.Expression {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Info("Storage.GetExpressions: Выдаем список выражений", "s.expressionList", s.expressionList)
	return s.expressionList
}

func (s *Storage) AddTask(task models.Task) int {
	slog.Info("Storage.AddTask: Получена задача", "task", task)
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = int(atomic.AddInt64(&s.taskID, 1))
	slog.Info("Storage.AddTask: Задаче назначен ID", "task", task)
	s.taskList = append(s.taskList, task)
	slog.Info("Storage.AddTask: Задача добавлена в список", "task", task, "s.taskList", s.taskList)
	return task.ID
}

func (s *Storage) GetPendingTask() *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.taskList {
		if s.taskList[i].Status == "pending" {
			slog.Info("Storage.GetPendingTask: Найдена задача", "s.taskList[i]", s.taskList[i])
			s.taskList[i].Status = "in progress"
			slog.Info("Storage.GetPendingTask: Задаче установлен статус в работе", "s.taskList[i]", s.taskList[i])
			return &s.taskList[i]
		}
	}
	return nil
}

func (s *Storage) DeleteTask(task_id int) {
	slog.Info("Storage.DeleteTask: Получена задача", "task_id", task_id)
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.taskList {
		if s.taskList[i].ID == task_id {
			slog.Info("Storage.DeleteTask: Задача найдена", "s.taskList[i]", s.taskList[i])
			s.taskList[i] = s.taskList[len(s.taskList)-1]
			s.taskList = s.taskList[:len(s.taskList)-1]
			slog.Info("Storage.DeleteTask: Задача удалена", "s.taskList", s.taskList)
			break
		}
	}
}

// Удаление всех задач, связанных с выражением
func (s *Storage) DeleteTaskByExpressionID(expression_id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.taskList {
		if s.taskList[i].ExpressionID == expression_id {
			slog.Info("Storage.DeleteTaskByExpressionID: Задача найдена", "s.taskList[i]", s.taskList[i])
			s.taskList[i] = s.taskList[len(s.taskList)-1]
			s.taskList = s.taskList[:len(s.taskList)-1]
			slog.Info("Storage.DeleteTaskByExpressionID: Задача удалена", "s.taskList", s.taskList)
			break
		}
	}
}

func (s *Storage) FindTaskByID(task_id int) *models.Task {
	slog.Info("Storage.FindTaskByID: Получена задача", "task_id", task_id)
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.taskList {
		if s.taskList[i].ID == task_id {
			slog.Info("Storage.FindTaskByID: Задача найдена", "s.taskList[i]", s.taskList[i])
			return &s.taskList[i]
		}
	}
	return nil
}

func (s *Storage) FindExpressionByID(expression_id int) *models.Expression {
	slog.Info("Storage.FindExpressionByID: Получено выражение", "expression_id", expression_id)
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.expressionList {
		if s.expressionList[i].ID == expression_id {
			slog.Info("Storage.FindExpressionByID: Выражение найдено", "s.expressionList[i]", s.expressionList[i])
			return &s.expressionList[i]
		}
	}
	return nil
}
