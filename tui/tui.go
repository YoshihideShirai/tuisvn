package tui

import (
	"strings"

	"github.com/rivo/tview"
)

type Tui struct {
	app               *tview.Application
	pages             *tview.Pages
	screen            map[string]map[string]*TuiScreen
	history_screen    map[int][]string
	root              *tview.Grid
	config            Config
	svnworker_limiter chan struct{}
}

type ConfigRepos struct {
	Url string `yaml:"url"`
}

type Config struct {
	Repos map[string]ConfigRepos `yaml:"repos"`
}

type TuiScreen struct {
	prim *tview.Grid
}

func New(repos_url string) *Tui {
	t := Tui{
		app:    tview.NewApplication(),
		screen: make(map[string]map[string]*TuiScreen),
	}
	svninfo := t.SvnInfo(repos_url)
	t.config.Repos = make(map[string]ConfigRepos)
	var config ConfigRepos
	config.Url = svninfo.Entry.Repository.Root

	path := strings.TrimPrefix(svninfo.Entry.Url, config.Url)
	if svninfo.Entry.Kind != "dir" {
		path_list := strings.Split(path, "/")
		path = strings.Join(path_list[:len(path_list)-1], "/")
	}
	path += "/"

	t.config.Repos[".arg"] = config
	t.screen[".arg"] = make(map[string]*TuiScreen)
	t.history_screen = make(map[int][]string)
	t.CreateApp(".arg", "tree:"+path)
	return &t
}

func NewRoot() *Tui {
	t := Tui{
		app:    tview.NewApplication(),
		screen: make(map[string]map[string]*TuiScreen),
	}
	t.screen[".root"] = make(map[string]*TuiScreen)
	t.history_screen = make(map[int][]string)
	t.CreateApp(".root", "main")
	return &t
}

func (t *Tui) Run() error {
	return t.app.Run()
}

func (t *Tui) CreateApp(repos string, screen string) {
	t.SvnWorkerInit()

	titlebar := TuiTitleBar("TuiSVN")
	grid := tview.NewGrid().
		SetRows(1, 0).
		SetBorders(false).
		AddItem(titlebar, 0, 0, 1, 3, 0, 0, false)
	t.root = grid
	t.app.SetRoot(grid, true)
	t.ChangeScreen(repos, screen)
}

func (t *Tui) ChangeScreen(repos string, screen string) {
	history_entry := []string{repos, screen}
	length := len(t.history_screen)
	t.history_screen[length] = history_entry
	t.changeScreenImpl(repos, screen)
}

func (t *Tui) BackScreen() {
	length := len(t.history_screen)
	if length < 2 {
		t.app.Stop()
		return
	}
	backscreen := t.history_screen[length-2]
	repos := backscreen[0]
	screen := backscreen[1]
	delete(t.history_screen, length-1)
	t.changeScreenImpl(repos, screen)
}

func (t *Tui) changeScreenImpl(repos string, screen string) {
	if t.screen[repos] == nil {
		t.screen[repos] = make(map[string]*TuiScreen)
	}
	if t.screen[repos][screen] == nil {
		if repos == ".root" && screen == "main" {
			t.NewTuiRoot()
		} else if strings.HasPrefix(screen, "tree:") {
			path := strings.TrimPrefix(screen, "tree:")
			t.NewTuiTree(repos, path)
		} else if strings.HasPrefix(screen, "log:") {
			path := strings.TrimPrefix(screen, "log:")
			t.NewTuiLog(repos, path)
		} else if strings.HasPrefix(screen, "rev:") {
			params := strings.Split(screen, ":")
			t.NewTuiRev(repos, params[1], params[2])
		} else if strings.HasPrefix(screen, "diff:") {
			params := strings.Split(screen, ":")
			t.NewTuiDiff(repos, params[1], params[2])
		} else {
			t.TuiPanic("bug:" + screen)
		}
	}
	t.root.
		AddItem(t.screen[repos][screen].prim, 1, 0, 1, 3, 0, 0, false)
	t.app.SetFocus(t.screen[repos][screen].prim)
}
