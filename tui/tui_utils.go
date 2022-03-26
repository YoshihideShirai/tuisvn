package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TuiTitleBar(text string) tview.Primitive {
	v := tview.NewTextView().
		SetTextAlign(tview.AlignLeft)
	v.SetBackgroundColor(tcell.ColorBlueViolet)
	v.SetText(text)
	return v
}

func TuiStatusBar(text string) tview.Primitive {
	v := tview.NewTextView().
		SetTextAlign(tview.AlignLeft)
	v.SetBackgroundColor(tcell.ColorBlue)
	v.SetText(text)
	return v
}
