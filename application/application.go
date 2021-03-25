package application

import (
	"sync"

	"github.com/corioders/gokit/log"
)

type Application struct {
	logger log.Logger

	stopHandlers []stopHandler
	stopFuncsMu  sync.Mutex
}

func New(logger log.Logger) *Application {
	return &Application{
		logger: logger,

		stopHandlers: make([]stopHandler, 0),
		stopFuncsMu:  sync.Mutex{},
	}
}

func (a *Application) GetLogger() log.Logger {
	return a.logger
}
