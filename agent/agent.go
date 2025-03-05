package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type task struct {
	ID             int           `json:"id"`
	Arg1           float64       `json:"arg1"`
	Arg2           float64       `json:"arg2"`
	Operation      string        `json:"operation"`
	OperationTime  time.Duration `json:"operation_time"`
}

type solvedTask struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}

func solveTask(t task) solvedTask {
	solved := solvedTask{ID: t.ID}

	time.Sleep(t.OperationTime)

	switch t.Operation {
	case "+":
		solved.Result = t.Arg1 + t.Arg2
	case "-":
		solved.Result = t.Arg1 - t.Arg2
	case "*":
		solved.Result = t.Arg1 * t.Arg2
	case "/":
		if t.Arg2 == 0 {
			log.Printf("Ошибка: деление на 0 в задаче ID %d\n", t.ID)
			solved.Result = 0
		} else {
			solved.Result = t.Arg1 / t.Arg2
		}
	default:
		log.Printf("Ошибка: неизвестная операция %s в задаче ID %d\n", t.Operation, t.ID)
	}

	return solved
}

func worker(tasks <-chan task, results chan<- solvedTask, wg *sync.WaitGroup) {
	defer wg.Done()


	for t := range tasks {
		timer := time.NewTimer(t.OperationTime)
		<-timer.C
		solved := solveTask(t)
		results <- solved
	}
}

func RunAgent() {
	// ------------------- Берем из env разные переменные -------------------
	taskPort, exists := os.LookupEnv("PORT")
	if !exists {
		taskPort = "8080"
	}
	taskURN, exists := os.LookupEnv("TASK_URN")
	if !exists {
		taskURN = "/internal/task"
	}
	taskURL := "http://localhost"
	taskURI := taskURL+":"+taskPort+taskURN
	workerCountStr, exists := os.LookupEnv("COMPUTING_POWER")
	if !exists {
		workerCountStr = "10"
	}
	var workerCount int
	if num, err := strconv.Atoi(workerCountStr); err != nil {
		workerCount = num
	} else {
		workerCount = 10
	}
	// ----------------------------------------------------------------------
	inputCh := make(chan task, workerCount)
	outputCh := make(chan solvedTask, workerCount)
	var wg sync.WaitGroup

	// эта горутина постоянно просит задачи
	go func() {
		defer close(inputCh)
		for {
			resp, err := http.Get(taskURI)
			if resp == nil {
				log.Print("Сервер не отвечает")
				time.Sleep(time.Second)
				continue
			}
			if resp.StatusCode == http.StatusNotFound {
				log.Print("Задач нет")
				time.Sleep(time.Second)
				continue
			}
			if err != nil {
				log.Printf("Ошибка при получении задачи: %v\n", err)
				continue
			}
			log.Printf("Получен ответ %d от сервера", resp.StatusCode)
			var t task
			if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
				log.Printf("Ошибка при декодировании JSON: %v\n", err)
			}
			log.Printf("Получена задача %v", t)
			resp.Body.Close()

			inputCh <- t
		}
	}()

	// Запускаем воркеров
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(inputCh, outputCh, &wg)
	}

	// горутина, которая отправляет решения
	go func() {
		for res := range outputCh {
			log.Printf("Отправляем решение %v", res)
			data, err := json.Marshal(res)
			if err != nil {
				log.Printf("Ошибка при маршалинге JSON: %v\n", err)
				continue
			}

			resp, err := http.Post(taskURI, "application/json", bytes.NewReader(data))
			if err != nil {
				log.Printf("Ошибка при отправке результата: %v\n", err)
				continue
			}
			resp.Body.Close()
		}
	}()

	wg.Wait()
	close(outputCh)
}