package main

import (
	"os"

	"github.com/YoshihideShirai/tuisvn/tui"
)

func main() {
	repos_url := "."

	if len(os.Args) >= 2 {
		repos_url = os.Args[1]
	}

	t := tui.New(repos_url)
	if err := t.Run(); err != nil {
		panic(err)
	}
}
