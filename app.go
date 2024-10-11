package main

import (
	"context"

	runtime "github.com/wailsapp/wails/v2/pkg/runtime"

	SVMiddleware "sqlview/middleware"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DatabaseButtonClicked(path string) {
	SVMiddleware.GetDatabase(path)
}

func (a *App) FileDialog() (string, error) {
	opts := runtime.OpenDialogOptions{
		DefaultDirectory: ".",
	}

	res, err := runtime.OpenFileDialog(a.ctx, opts)
	if err != nil {
		return "", err
	}
	if res == "" {
		return "", nil
	}

	SVMiddleware.GetDatabase(res)
	return res, nil
}

func (a *App) Query(query string) (string, error) {
	ret, err := SVMiddleware.ProcessQuery(query)
	if err != nil {
		return "", err
	}
	return ret, nil
}
