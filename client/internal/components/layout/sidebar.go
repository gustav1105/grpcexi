package layout

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sidebar struct {
	*tview.Table
	selectedFunc func(index int, text string)
	items        []string
}

func NewSidebar(title string) *Sidebar {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)
	table.SetBorder(true).
		SetTitle(title)

	s := &Sidebar{Table: table}
	s.initKeybindings()
	return s
}

func (s *Sidebar) SetItems(items []string, selectedFunc func(index int, mainText string)) {
	s.Clear()
	s.items = items
	s.selectedFunc = selectedFunc

	for i, item := range items {
		cell := tview.NewTableCell(item).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetSelectable(true)
		s.SetCell(i, 0, cell)
	}

	if len(items) > 0 {
		s.Select(0, 0)
	}
}

func (s *Sidebar) initKeybindings() {
	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := s.GetSelection()
		switch event.Key() {
		case tcell.KeyUp:
			if row > 0 {
				s.Select(row-1, 0)
			}
			return nil
		case tcell.KeyDown:
			if row < len(s.items)-1 {
				s.Select(row+1, 0)
			}
			return nil
		case tcell.KeyEnter:
			if s.selectedFunc != nil && row >= 0 && row < len(s.items) {
				s.selectedFunc(row, s.items[row])
			}
			return nil
		}
		return event
	})
}

