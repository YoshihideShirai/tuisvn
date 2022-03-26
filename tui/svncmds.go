package tui

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
)

type SvnLogXml struct {
	Logentry []struct {
		Revision string `xml:"revision,attr"`
		Author   string `xml:"author"`
		Date     string `xml:"date"`
		Msg      string `xml:"msg"`
		Path     []struct {
			PropMods     string `xml:"prop-mods,attr"`
			TextMods     string `xml:"text-mods,attr"`
			Kind         string `xml:"kind,attr"`
			Action       string `xml:"action,attr"`
			CopyFromPath string `xml:"copyfrom-path,attr"`
			CopyFromRev  string `xml:"copyfrom-rev,attr"`
			Path         string `xml:",chardata"`
		} `xml:"paths>path"`
	} `xml:"logentry"`
}

type SvnInfoXml struct {
	Entry struct {
		Kind       string `xml:"kind,attr"`
		Path       string `xml:"path,attr"`
		Revision   string `xml:"revision,attr"`
		Url        string `xml:"url"`
		Repository struct {
			Root string `xml:"root"`
			Uuid string `xml:"uuid"`
		} `xml:"repository"`
	} `xml:"entry"`
}

func (t *Tui) SvnWorkerInit() {
	t.svnworker_limiter = make(chan struct{}, 10)
}

func (t *Tui) SvnLs(repos string, path string) string {
	url := t.config.Repos[repos].Url + path
	cmd := exec.Command("svn", "ls", url)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.TuiPanic(stderr.String())
	}
	return stdout.String()
}

func (t *Tui) SvnDiff(repos string, path string, rev string) string {
	url := t.config.Repos[repos].Url + path
	rev = strings.TrimPrefix(rev, "r")
	rev_opt := "-c" + rev
	cmd := exec.Command("svn", "diff", rev_opt, url)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.TuiPanic(stderr.String())
	}
	return stdout.String()
}

func (t *Tui) SvnLogSummary(repos string, path string) *SvnLogXml {
	t.svnworker_limiter <- struct{}{}
	res := t.SvnLog(repos, path, "HEAD", "1", 1)
	<-t.svnworker_limiter
	return res
}

func (t *Tui) SvnLog(repos string, path string, fromrev string, torev string, count int) *SvnLogXml {
	svnlog := new(SvnLogXml)
	url := t.config.Repos[repos].Url + path
	cmd := exec.Command("svn", "log",
		"-l", fmt.Sprintf("%d", count),
		"-r", fmt.Sprintf("%s:%s", fromrev, torev),
		"-v", "--xml", url)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.TuiPanic(stderr.String())
	}
	if err := xml.Unmarshal([]byte(stdout.String()), svnlog); err != nil {
		t.TuiPanic(err.Error())
	}
	return svnlog
}

func (t *Tui) SvnInfo(url string) *SvnInfoXml {
	svnlog := new(SvnInfoXml)
	cmd := exec.Command("svn", "info",
		"--xml", url)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.TuiPanic(stderr.String())
	}
	if err := xml.Unmarshal([]byte(stdout.String()), svnlog); err != nil {
		t.TuiPanic(err.Error())
	}
	return svnlog
}
