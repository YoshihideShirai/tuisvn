package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) NewTuiRoot() {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}
	statusbar := TuiStatusBar("[.root]")
	main := tview.NewTable().SetSelectable(true, false)

	idx := 0
	for i, v := range t.config.Repos {
		main.SetCell(idx, 0, tview.NewTableCell(i))
		main.SetCell(idx, 1, tview.NewTableCell(v.Url).SetExpansion(1))
		idx++
	}

	s.prim.
		SetRows(0, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 1, 0, 1, 3, 0, 0, false)

	s.prim.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := main.GetSelection()
			cell := main.GetCell(row, 0)
			repos := cell.Text
			t.ChangeScreen(repos, "tree:/")
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'k':
				row, _ := main.GetSelection()
				row--
				main.Select(row, 0)
				return nil
			case 'j':
				row, _ := main.GetSelection()
				if row < main.GetRowCount()-1 {
					row++
				}
				main.Select(row, 0)
				return nil
			case 'q':
				t.BackScreen()
				return nil
			}
		}
		return event
	})
	t.screen[".root"]["main"] = &s
}
