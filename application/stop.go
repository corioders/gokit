package application

import "fmt"

type StopRegistrar interface {
	RegisterOnStop(name string, fn stopFunc)
}

type stopFunc func() error
type stopHandler struct {
	fn   stopFunc
	name string
}

// RegisterOnStop registers function that should be executed when application is stopping.
func (a *Application) RegisterOnStop(name string, fn stopFunc) {
	a.onStopMu.Lock()
	defer a.onStopMu.Unlock()
	a.onStop = append(a.onStop, stopHandler{name: name, fn: fn})
}

func (a *Application) Stop() error {
	a.logger.Info("Stopping...")

	a.onStopMu.Lock()
	defer a.onStopMu.Unlock()
	for _, handler := range a.onStop {
		a.logger.Info(fmt.Sprintf("Stopping %s...", handler.name))

		err := handler.fn()
		if err != nil {
			return err
		}

		a.logger.Info(fmt.Sprintf("Stopped successfully %s...", handler.name))
	}

	a.logger.Info("Stopped successfully...")
	return nil
}
