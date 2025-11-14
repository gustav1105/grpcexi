package modal

import "github.com/rivo/tview"

type Quit struct {
	*tview.Modal
}

func NewQuit(app *tview.Application, pages *tview.Pages) *Quit {
	quit := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			} else if buttonLabel == "Cancel" {
				pages.SwitchToPage("home")
			}
		})

	return &Quit{
		Modal: quit,
	}
}
