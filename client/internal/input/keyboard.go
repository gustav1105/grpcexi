package input

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Handler struct {
	app   *tview.Application
	pages *tview.Pages
}

func NewHandler(app *tview.Application, pages *tview.Pages) *Handler {
	return &Handler{app: app, pages: pages}
}

func (h *Handler) CaptureKeys() {
	h.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			h.pages.SwitchToPage("quit")
			return nil
		}

		return event
	})
}
