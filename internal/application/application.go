package application

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/romanSPB15/Calculator_Service/internal/web"
	"github.com/romanSPB15/Calculator_Service/pckg/dir"
	"github.com/romanSPB15/Calculator_Service/pckg/rpn"
)

const (
	WaitStatus        = "Wait"
	OKStatus          = "OK"
	CalculationStatus = "Calculation"
)

// Выражение
type Expression struct {
	Data   string  `json:"data"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// Выражение с ID
type ExpressionWithID struct {
	ID IDExpression `json:"id"`
	Expression
}

// ID выражения
type IDExpression = uint32

type GetExpressionHandlerResult struct {
	Expression ExpressionWithID `json:"expression"`
}

type AddHandlerResult struct {
	ID uint32 `json:"id"`
}

type GetExpressionsHandlerResult struct {
	Expressions []ExpressionWithID `json:"expressions"`
}

type GetTaskHandlerResult struct {
	Task rpn.TaskID `json:"task"`
}

type AgentResult struct {
	ID     rpn.IDTask `json:"id"`
	Result float64    `json:"result"`
}

// Выражения
var Expressions = make(map[IDExpression]*Expression)

// Задачи
var Tasks = rpn.NewConcurrentTaskMap()

// Приложение
type Application struct {
	// Агент
	Config       *config
	Agent        http.Client
	NumGoroutine int
	Router       *mux.Router
}

func New() *Application {
	return &Application{
		Router: mux.NewRouter(),
		Config: newConfig(),
	}
}

// Запуск всей системы
func (app *Application) RunServer() {
	rpn.InitEnv(dir.EnvFile()) // Иницилизация переменных из среды
	/* ListenAndServe() закончится только с ошибкой, как и runAgent() */
	startServer := make(chan struct{}, 1) // Канал запуска оркестратора
	go func() {
		startServer <- struct{}{}
		if app.Config.Debug {
			log.Println("Orkestrator Runned")
		}
		err := http.ListenAndServe(":8080", nil)
		panic(err)
	}()
	// Создаём новый mux.Router
	/* Инициализация обработчиков роутера */
	app.Router.HandleFunc("/api/v1/calculate", app.AddExpressionHandler)
	app.Router.HandleFunc("/api/v1/expressions/{id}", app.GetExpressionHandler)
	app.Router.HandleFunc("/api/v1/expressions", app.GetExpressionsHandler)
	app.Router.HandleFunc("/api/v1/internal/task", app.TaskHandler)
	if app.Config.Web {
		web.HandleToRouter(app.Router)
	}
	http.Handle("/", app.Router)
	<-startServer // Ждём, когда запустится оркестратор
	panic(app.runAgent())
}
