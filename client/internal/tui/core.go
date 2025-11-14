package tui

import (
	"grpcexi/client/internal/components/layout"
	"grpcexi/client/internal/handlers"
	"grpcexi/client/internal/input"
	"grpcexi/client/internal/page"
	"grpcexi/client/internal/page/modal"
	"grpcexi/client/internal/services"
	pb "grpcexi/protos"

	"github.com/rivo/tview"
)

type App struct {
	app    *tview.Application
	pages  *tview.Pages
	client pb.ContactServiceClient
}

func NewApp(client pb.ContactServiceClient) *App {
	app := tview.NewApplication()
	pages := tview.NewPages()

	service := services.NewContactService(client)
	handler := handlers.NewContactHandler(app, service)

	home := page.NewHome(app, handler)
	quitModal := modal.NewQuit(app, pages)

	menuBar := layout.NewMenuBar()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(home, 0, 1, true).
		AddItem(menuBar, 1, 0, false)

	pages.AddPage("home", flex, true, true)
	pages.AddPage("quit", quitModal, true, false)

	inputHandler := input.NewHandler(app, pages)
	inputHandler.CaptureKeys()

	return &App{
		app:    app,
		pages:  pages,
		client: client,
	}
}

func (a *App) GetRoot() tview.Primitive {
	return a.pages
}

func (a *App) Run() error {
	return a.app.SetRoot(a.pages, true).Run()
}
