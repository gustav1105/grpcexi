package tui

import "github.com/rivo/tview"

type TUI interface {
	GetRoot() tview.Primitive
	Run() error
}
