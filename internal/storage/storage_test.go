package storage

import (
	"sync"
	"testing"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage_AddExpression(t *testing.T) {
	s := NewStorage()
	expr := models.Expression{}
	id := s.AddExpression(expr)

	assert.Equal(t, 1, id, "ID of the first expression should be 1")
	assert.Len(t, s.GetExpressions(), 1, "There should be one expression in the storage")
}

func TestStorage_AddTask(t *testing.T) {
	s := NewStorage()
	task := models.Task{}
	id := s.AddTask(task)

	assert.Equal(t, 1, id, "ID of the first task should be 1")
	assert.Len(t, s.taskList, 1, "There should be one task in the storage")
}

func TestStorage_GetPendingTask(t *testing.T) {
	s := NewStorage()
	task1 := models.Task{Status: "pending"}
	task2 := models.Task{Status: "in progress"}
	s.AddTask(task1)
	s.AddTask(task2)

	pendingTask := s.GetPendingTask()

	assert.NotNil(t, pendingTask, "Expected task should be found")
	assert.Equal(t, "in progress", pendingTask.Status, "Task status should change to 'in progress'")
}

func TestStorage_DeleteTask(t *testing.T) {
	s := NewStorage()
	task := models.Task{}
	id := s.AddTask(task)

	s.DeleteTask(id)
	assert.Nil(t, s.FindTaskByID(id), "Task should be deleted")
}

func TestStorage_FindTaskByID(t *testing.T) {
	s := NewStorage()
	task := models.Task{}
	id := s.AddTask(task)

	foundTask := s.FindTaskByID(id)
	assert.NotNil(t, foundTask, "Task should be found")
	assert.Equal(t, id, foundTask.ID, "ID of the found task should match")
}

func TestStorage_FindExpressionByID(t *testing.T) {
	s := NewStorage()
	expr := models.Expression{}
	id := s.AddExpression(expr)

	foundExpr := s.FindExpressionByID(id)
	assert.NotNil(t, foundExpr, "Expression should be found")
	assert.Equal(t, id, foundExpr.ID, "ID of the found expression should match")
}

func TestStorage_Concurrency(t *testing.T) {
	s := NewStorage()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			s.AddTask(models.Task{})
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, 100, len(s.taskList), "All 100 tasks should be added correctly")
}
