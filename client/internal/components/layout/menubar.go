package layout

import "github.com/rivo/tview"

func NewMenuBar() tview.Primitive {
	page1 := tview.NewTextView().
		SetText("Home").
		SetTextAlign(tview.AlignLeft)

	menuBar := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(page1, 0, 1, false)

	return menuBar
}
