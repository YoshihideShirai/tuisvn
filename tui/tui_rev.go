package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) NewTuiRev(repos string, path string, rev string) {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}
	statusbar := TuiStatusBar(fmt.Sprintf("[%s]rev:%s:%s", repos, path, rev))
	main := tview.NewTable().SetSelectable(true, false)
	idx := 0
	revlog := t.SvnLog(repos, path, rev, rev, 1).Logentry[0]
	main.SetCell(idx, 0, tview.NewTableCell(fmt.Sprintf(
		"[yellow]Revision : %s[white]", tview.Escape(revlog.Revision))).SetExpansion(1))
	idx++
	main.SetCell(idx, 0, tview.NewTableCell(fmt.Sprintf(
		"[blue]Date     : %s[white]", tview.Escape(revlog.Date))).SetExpansion(1))
	idx++
	main.SetCell(idx, 0, tview.NewTableCell(fmt.Sprintf(
		"[pink]Author   : %s[white]", tview.Escape(revlog.Author))).SetExpansion(1))
	idx++
	main.SetCell(idx, 0, tview.NewTableCell("").SetExpansion(1))
	idx++

	for _, v := range strings.Split(revlog.Msg, "\n") {
		main.SetCell(idx, 0, tview.NewTableCell(fmt.Sprintf(
			"    %s", tview.Escape(v))).SetExpansion(1))
		idx++
	}
	main.SetCell(idx, 0, tview.NewTableCell("").SetExpansion(1))
	idx++

	main.SetCell(idx, 0, tview.NewTableCell("[blue]Changed paths:[white]").SetExpansion(1))
	idx++

	path_idx_head := idx
	for _, v := range revlog.Path {
		color := "white"
		switch v.Action {
		case "A":
			color = "green"
		case "M":
			color = "yellow"
		case "D":
			color = "red"
		default:
		}
		celltext := fmt.Sprintf("    %s %s", tview.Escape(v.Action), tview.Escape(v.Path))
		if v.CopyFromPath != "" {
			celltext += fmt.Sprintf(" (from %s:%s)", v.CopyFromPath, v.CopyFromRev)
		}
		main.SetCell(idx, 0,
			tview.NewTableCell(fmt.Sprintf("[%s]%s[white]",
				color, celltext)).SetExpansion(1))
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
			path_idx := row - path_idx_head
			if path_idx >= 0 {
				diff_path := revlog.Path[path_idx]
				change_screen := "diff:" + diff_path.Path + ":" + rev
				t.ChangeScreen(repos, change_screen)
			}
			return nil
		case tcell.KeyDown:
			row, _ := main.GetSelection()
			if row < main.GetRowCount()-1 {
				row++
			}
			main.Select(row, 0)
			return nil
		case tcell.KeyUp:
			row, _ := main.GetSelection()
			row--
			main.Select(row, 0)
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
	t.screen[repos]["rev:"+path+":"+rev] = &s
}
