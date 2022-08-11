package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SummaryLogRequest struct {
	repos string
	path  string
	idx   int
	file  string
}

func (t *Tui) TuiTreeUpdateSummary(ch chan SummaryLogRequest, repos string, path string, idx int, file string) {
	var req SummaryLogRequest
	req.repos = repos
	req.path = path
	req.idx = idx
	req.file = file
	ch <- req
}

func (t *Tui) TuiTreeUpdateSummaryWorker(ch chan SummaryLogRequest, table *tview.Table) {
	for {
		req, more := <-ch
		if more == false {
			return
		}
		repos := req.repos
		path := req.path
		file := req.file
		idx := req.idx

		svnlog := t.SvnLogSummary(repos, path+file)
		table.SetCell(idx, 1,
			tview.NewTableCell(fmt.Sprintf("[yellow]%s[white]", tview.Escape(
				"r"+svnlog.Logentry[0].Revision))))
		table.SetCell(idx, 2,
			tview.NewTableCell(fmt.Sprintf("[green]%s[white]", tview.Escape(
				strings.Split(svnlog.Logentry[0].Msg, "\n")[0]))).
				SetExpansion(1))
		t.app.Draw()
	}
}

func (t *Tui) TuiTreeUpdate(repos string, path string, table *tview.Table) {
	o := t.SvnLs(repos, path)
	idx := 0
	table.Select(0, 0)
	table.Clear()

	table.SetCell(idx, 0, tview.NewTableCell("."))
	table.SetCell(idx, 1, tview.NewTableCell(""))
	table.SetCell(idx, 2, tview.NewTableCell("").SetExpansion(1))
	idx++

	if path != "/" {
		table.SetCell(idx, 0, tview.NewTableCell(".."))
		table.SetCell(idx, 1, tview.NewTableCell(""))
		table.SetCell(idx, 2, tview.NewTableCell("").SetExpansion(1))
		idx++
	}

	worker_ch := make(chan SummaryLogRequest, 256)
	for i := 0; i < 10; i++ {
		go t.TuiTreeUpdateSummaryWorker(worker_ch, table)
	}
	filedirs_list := strings.Split(o, "\n")
	sort.Slice(filedirs_list, func(i, j int) bool {
		i_is_dir := strings.HasSuffix(filedirs_list[i], "/")
		j_is_dir := strings.HasSuffix(filedirs_list[j], "/")
		if i_is_dir != j_is_dir {
			return i_is_dir
		}
		return filedirs_list[i] < filedirs_list[j]
	})

	t.TuiTreeUpdateSummary(worker_ch, repos, "/", 0, "")

	for _, v := range filedirs_list {
		if len(v) > 0 {
			table.SetCell(idx, 0,
				tview.NewTableCell(fmt.Sprintf("%s", tview.Escape(v))))
			table.SetCell(idx, 1, tview.NewTableCell(""))
			table.SetCell(idx, 2, tview.NewTableCell("").SetExpansion(1))
			t.TuiTreeUpdateSummary(worker_ch, repos, path, idx, v)
			idx++
		}
	}
	close(worker_ch)
}

func (t *Tui) NewTuiTree(repos string, path string) {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}
	statusbar := TuiStatusBar(fmt.Sprintf("[%s]tree:%s", repos, path))
	main := tview.NewTable().SetSelectable(true, false)
	s.prim.
		SetRows(0, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 1, 0, 1, 3, 0, 0, false)

	t.TuiTreeUpdate(repos, path, main)

	s.prim.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := main.GetSelection()
			fi := main.GetCell(row, 0).Text
			if fi == "." {
				return nil
			}
			if strings.HasSuffix(fi, "/") || fi == ".." {
				var changedpath string
				if fi == ".." {
					path_list := strings.Split(path, "/")
					changedpath = strings.Join(path_list[:len(path_list)-2], "/") + "/"
				} else {
					changedpath = path + fi
				}
				t.ChangeScreen(repos, "tree:"+changedpath)
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
			case 'l':
				row, _ := main.GetSelection()
				fi := main.GetCell(row, 0).Text
				if fi == ".." {
					return nil
				}
				changedpath := path + fi
				t.ChangeScreen(repos, "log:"+changedpath)
				return nil
			case 'q':
				t.BackScreen()
				return nil
			}
		}
		return event
	})
	t.screen[repos]["tree:"+path] = &s
}
