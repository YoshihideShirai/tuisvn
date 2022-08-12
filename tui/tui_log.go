package tui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) TuiLogUpdateWorker(repos string, path string, table *tview.Table) {
	read_rev := "HEAD"
	idx := 0
	log_length := 10
	for {
		first_rev := read_rev
		if first_rev != "HEAD" {
			read_rev_int, _ := strconv.Atoi(read_rev)
			first_rev_int := read_rev_int - 1
			if first_rev_int < 1 {
				first_rev_int = 1
			}
			first_rev = strconv.Itoa(first_rev_int)
		}
		res := t.SvnLog(repos, path, first_rev, "1", log_length)
		for _, v := range res.Logentry {
			table.SetCell(idx, 0,
				tview.NewTableCell(fmt.Sprintf("[blue]%s[white]", tview.Escape(
					v.Date))))
			table.SetCell(idx, 1,
				tview.NewTableCell(fmt.Sprintf("[green]%s[white]", tview.Escape(
					v.Author))))
			table.SetCell(idx, 2,
				tview.NewTableCell(fmt.Sprintf("[yellow]%s[white]", tview.Escape(
					"r"+v.Revision))))
			table.SetCell(idx, 3,
				tview.NewTableCell(fmt.Sprintf("[white]%s[white]", tview.Escape(
					strings.Split(v.Msg, "\n")[0]))).
					SetExpansion(1))
			idx++
			read_rev = v.Revision
		}
		t.app.Draw()
		if len(res.Logentry) < log_length {
			break
		}
	}
}

func (t *Tui) NewTuiLog(repos string, path string) {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}
	statusbar := TuiStatusBar(fmt.Sprintf("[%s]log:%s", repos, path))
	main := tview.NewTable().SetSelectable(true, false)
	main.SetCell(0, 0, tview.NewTableCell(""))
	main.SetCell(0, 1, tview.NewTableCell(""))
	main.SetCell(0, 2, tview.NewTableCell(""))
	main.SetCell(0, 3, tview.NewTableCell("").SetExpansion(1))

	s.prim.
		SetRows(0, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 1, 0, 1, 3, 0, 0, false)

	go t.TuiLogUpdateWorker(repos, path, main)

	s.prim.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := main.GetSelection()
			rev := main.GetCell(row, 2).Text
			re1 := regexp.MustCompile(`\[[A-Za-z]+\]`)
			rev = re1.ReplaceAllString(rev, "")
			t.ChangeScreen(repos, "rev:"+path+":"+rev)
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
	t.screen[repos]["log:"+path] = &s
}
