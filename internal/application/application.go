package application

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/Yshariale/FinalTaskFirstSprint/internal/config"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/services/expression"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/storage"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/transport/handlers"
	"github.com/Yshariale/FinalTaskFirstSprint/internal/transport/middlewares"
)

func setUpLogger(logFile *os.File) error {
	var logger = slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)
	return nil
}

type Application struct {
	config *config.Config
}

func New() *Application {
	return &Application{
		config: config.ConfigFromEnv(),
	}
}

func (a *Application) RunServer() error {
	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error while opening log file", "error", err)
	}
	defer logFile.Close()
	err = setUpLogger(logFile)
	if err != nil {
		slog.Error("Error while setting up logger", "error", err)
	}

	// Создаем хранилище, которое будет передаваться вглубь приложения по ссылке,
	// то есть все сервисы будут работать с одним и тем же хранилищем
	storage := storage.NewStorage()

	// А вот и сервис по работе с выражениями. Он используется в хендлерах для обработки запросов
	expressionService := expression.NewExpressionService(storage, a.config.TimeConf)

	slog.Info("Starting server", "port", a.config.Addr)
	r := mux.NewRouter()
	r.Handle("/api/v1/calculate", middlewares.LoggingMiddleware(handlers.NewCalcHandler(expressionService))).Methods(http.MethodPost)
	r.Handle("/api/v1/expressions", middlewares.LoggingMiddleware(handlers.NewExpressionListHandler(expressionService))).Methods(http.MethodGet)
	r.Handle("/api/v1/expressions/{id:[0-9]+}", middlewares.LoggingMiddleware(handlers.NewExpressionHandler(expressionService))).Methods(http.MethodGet)
	r.Handle("/internal/task", middlewares.LoggingMiddleware(handlers.NewTaskHandler(expressionService)))
	http.Handle("/", r)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
