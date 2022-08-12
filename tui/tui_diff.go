package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) NewTuiDiff(repos string, path string, rev string) {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}

	diff_output := t.SvnDiff(repos, path, rev)

	statusbar := TuiStatusBar(fmt.Sprintf("[%s]diff:%s:%s", repos, path, rev))
	main := tview.NewTable().SetSelectable(true, false)

	for i, v := range strings.Split(diff_output, "\n") {
		color := "white"
		if strings.HasPrefix(v, "+") {
			color = "green"
		} else if strings.HasPrefix(v, "-") {
			color = "red"
		} else if strings.HasPrefix(v, "@") {
			color = "purple"
		}
		output := fmt.Sprintf("[%s]%s[white]", color, tview.Escape(v))
		main.SetCell(i, 0,
			tview.NewTableCell(output).SetExpansion(1))
	}

	s.prim.
		SetRows(0, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 1, 0, 1, 3, 0, 0, false)

	s.prim.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
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
		case tcell.KeyUp:
			row, _ := main.GetSelection()
			row--
			main.Select(row, 0)
			return nil
		case tcell.KeyDown:
			row, _ := main.GetSelection()
			if row < main.GetRowCount()-1 {
				row++
			}
			main.Select(row, 0)
			return nil
		}
		return event
	})
	t.screen[repos]["diff:"+path+":"+rev] = &s
}
