package application

import "fmt"

type StopHandler interface {
	StopFunc(name string, fn stopFunc)
}

type stopFunc func() error
type stopHandler struct {
	fn   stopFunc
	name string
}

// StopFunc registers function that should be executed when application is stopping
func (a *Application) StopFunc(name string, fn stopFunc) {
	a.stopFuncsMu.Lock()
	defer a.stopFuncsMu.Unlock()
	a.stopHandlers = append(a.stopHandlers, stopHandler{name: name, fn: fn})
}

func (a *Application) Stop() error {
	a.logger.Info("Stopping...")

	a.stopFuncsMu.Lock()
	defer a.stopFuncsMu.Unlock()
	for _, handler := range a.stopHandlers {
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
