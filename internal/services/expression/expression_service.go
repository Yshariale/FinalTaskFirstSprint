package expression

//
//
// Этот модуль содержит логику обработки выражений
// ExpressionService взаимодействует со списком заданий
// и выражений через хранилище Storage
//
// Здесь много много логов, но это все нужно. Не удалять же...
//

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/config"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/models"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/calculation"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/storage"
)

type ExpressionService struct {
	storage    *storage.Storage
	timeConfig config.TimeConfig
}

func NewExpressionService(s *storage.Storage, tc config.TimeConfig) *ExpressionService {
	return &ExpressionService{storage: s, timeConfig: tc}
}

// Обработчик входящего выражения.
// Он запускается один раз для каждого выражения
func (s *ExpressionService) ProcessExpression(expressionStr string) (int, error) {
	slog.Info("ExpressionService.ProcessExpression: Начало обработки выражения", "expression", expressionStr)
	// Первым делом переводим в постфиксную запись
	postfix, err := calculation.ToPostfix(expressionStr)
	if err != nil {
		slog.Error("ExpressionService.ProcessExpression: Ошибка при переводе в постфиксную запись")
		return 0, err
	}
	slog.Info("ExpressionService.ProcessExpression: Выражение переведено в постфиксную запись", "expression", postfix)

	// Формируем выражение и здесь же строим бинарное дерево
	newExpression := models.Expression{
		Status:     "processing",
		BinaryTree: calculation.BuildTree(postfix),
	}
	slog.Info("ExpressionService.ProcessExpression: Построено бинарное дерево", "BinaryTree", newExpression.BinaryTree)

	// Добавляем выражение в хранилище
	expressionID := s.storage.AddExpression(newExpression)
	slog.Info("ExpressionService.ProcessExpression: Выражение добавлено в хранилище", "id", expressionID)

	// Ищем вершины у которых дети это числа...
	spareNodes := newExpression.BinaryTree.FindSpareNodes()
	slog.Info("ExpressionService.ProcessExpression: Получен список свободных узлов", "spareNodes", spareNodes)
	for _, node := range spareNodes {
		// ..., и создаем для них задачи
		s.createTaskForSpareNode(node, expressionID)
	}
	slog.Info("ExpressionService.ProcessExpression: Конец обработки выражения", "id", expressionID)
	return expressionID, nil
}

// Создание задачи для свободного узла. Свободный - это узел, у которого оба ребенка - числа
func (s *ExpressionService) createTaskForSpareNode(node *calculation.TreeNode, expressionID int) {
	slog.Info("ExpressionService.createTaskForSpareNode: Начало создания задачи для узла", "node", node, "expressionID", expressionID)
	arg1, _ := strconv.ParseFloat(node.Left.Val, 64)
	arg2, _ := strconv.ParseFloat(node.Right.Val, 64)
	slog.Info("ExpressionService.createTaskForSpareNode: Извлечены аргументы для задачи", "arg1", arg1, "arg2", arg2)

	if arg2 == 0 && node.Val == "/" {
		// если делим на ноль, то закрываем выражение
		slog.Info("ExpressionService.createTaskForSpareNode: Деление на 0, закрываем выражение", "node", node)
		s.closeExpressionWithError(s.storage.FindExpressionByID(expressionID), "division by zero")
	}

	task := models.Task{
		ExpressionID:  expressionID,
		Status:        "pending",
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     node.Val,
		OperationTime: s.getOperationTime(node.Val),
	}
	slog.Info("ExpressionService.createTaskForSpareNode: Сформирована задача", "task", task)

	// Добавляем задачу в хранилище
	taskID := s.storage.AddTask(task)
	slog.Info("ExpressionService.createTaskForSpareNode: Задача добавлена в хранилище", "taskID", taskID)
	node.TaskID = taskID
	slog.Info("ExpressionService.createTaskForSpareNode: Вершине присвоена ID задачи", "node", node, "taskID", taskID)
}

func (s ExpressionService) getOperationTime(operation string) time.Duration {
	switch operation {
	case "+":
		return s.timeConfig.TimeAdd
	case "-":
		return s.timeConfig.TimeSub
	case "*":
		return s.timeConfig.TimeMul
	case "/":
		return s.timeConfig.TimeDiv
	default:
		return 0
	}
}

// Получение списка выражений из хранилища
func (s *ExpressionService) GetExpressions() []models.Expression {
	slog.Info("ExpressionService.GetExpressions: Выдаем список выражений")
	return s.storage.GetExpressions()
}

// Получение выражения по ID
func (s *ExpressionService) GetExpressionByID(id int) *models.Expression {
	return s.storage.FindExpressionByID(id)
}

// Этот метод раздает задачу, которая ждет отправки
func (s *ExpressionService) GetPendingTask() *models.Task {
	return s.storage.GetPendingTask()
}

// Обработка входящей задачи. Или по другому: запускается когда агент отправляет результат задачи
func (s *ExpressionService) ProcessIncomingTask(task_id int, result float64) {
	slog.Info("ExpressionService.ProcessIncomingTask: Начало обработки входящей задачи", "task_id", task_id, "result", result)
	task := s.storage.FindTaskByID(task_id)
	slog.Info("ExpressionService.ProcessIncomingTask: Найдена задача по task_id", "task_id", task_id, "task", task)
	task.Status = "done"
	slog.Info("ExpressionService.ProcessIncomingTask: Задаче установлен статус done", "task", task)
	expression := s.storage.FindExpressionByID(task.ExpressionID)
	slog.Info("ExpressionService.ProcessIncomingTask: Найдено выражение для задачи", "task.ExpressionID", task.ExpressionID, "expression", expression)
	s.storage.DeleteTask(task_id)
	slog.Info("ExpressionService.ProcessIncomingTask: Задача удалена", "task_id", task_id)
	// Здесь самое интересное. Когда пришел результат задачи мы заменяем вершину задачи на результат...
	parent_task_node, node := expression.BinaryTree.FindParentAndNodeByTaskID(task_id)
	slog.Info("ExpressionService.ProcessIncomingTask: Найден узел и родитель узла для задачи", "task_id", task_id, "node", node, "parent_task_node", parent_task_node)
	expression.BinaryTree.ReplaceNodeWithValue(node, result)
	slog.Info("ExpressionService.ProcessIncomingTask: Узел заменен на значение", "expression.BinaryTree", expression.BinaryTree)
	if parent_task_node == nil {
		// ... если у вершины нет родителя, то это значит, что это корень дерева и выражение решено
		slog.Info("ExpressionService.ProcessIncomingTask: У узла нет родителя, завершаем вычисление выражения")
		s.solveExpression(expression, result)
		return
	}
	// ... и проверяем, можно ли из родителя сделать задачу
	if parent_task_node.IsSpare() {
		slog.Info("ExpressionService.ProcessIncomingTask: Из родителя можно сделать задачу")
		s.createTaskForSpareNode(parent_task_node, expression.ID)
	}
	slog.Info("ExpressionService.ProcessIncomingTask: Из родителя нельзя сделать задачу")
}

func (s *ExpressionService) closeExpressionWithError(expression *models.Expression, errorMsg string) {
	expression.Status = "error " + errorMsg
	expression.BinaryTree = nil
	slog.Info("ExpressionService.CloseExpressionWithError: Закрытие выражения с ошибкой", "expression", expression, "error", errorMsg)
	s.storage.DeleteTaskByExpressionID(expression.ID)
}

func (s *ExpressionService) solveExpression(expression *models.Expression, result float64) {
	expression.Result = result
	expression.Status = "solve"
	expression.BinaryTree = nil
	slog.Info("ExpressionService.SolveExpression: Выражение решено", "expression", expression)
}
