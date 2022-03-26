package tui

import (
	"fmt"
	"os"
)

func (t *Tui) TuiPanic(v string) {
	t.app.Stop()
	if DEBUG {
		panic(v)
	} else {
		fmt.Fprintln(os.Stderr, v)
	}
	os.Exit(-1)
}
