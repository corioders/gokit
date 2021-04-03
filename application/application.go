package application

import (
	"sync"

	"github.com/corioders/gokit/log"
)

type Application struct {
	logger log.Logger

	onStop []stopHandler
	onStopMu  sync.Mutex
}

func New(logger log.Logger) *Application {
	return &Application{
		logger: logger,

		onStop: make([]stopHandler, 0),
		onStopMu:  sync.Mutex{},
	}
}

func (a *Application) GetLogger() log.Logger {
	return a.logger
}
