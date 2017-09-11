package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	lineBreak = "\n"

	fileStartPrefix = "-- # start # "
	fileEndPrefix   = "-- # end # "
)

var (
	colorOK       = color.New(color.BgGreen, color.FgWhite, color.Bold)
	colorFail     = color.New(color.BgRed, color.FgWhite, color.Bold)
	colorTestPath = color.New()
	colorMessage  = color.New(color.Bold)
	colorFailPath = color.New(color.Underline)
	colorFailChar = color.New(color.FgRed, color.Underline)
	colorText     = color.New()
)

type logger interface {
	ok()
	failSQL(message, text string, pos int)
	fail(string)
}

type stdLog struct {
	path string
}

func newStdLog(path string) *stdLog {
	return &stdLog{path}
}

func (l *stdLog) ok() {
	fmt.Println(colorOK.Sprintf("  OK  ") + " " + colorTestPath.Sprintf(l.path))
}

func (l *stdLog) fail(msg string) {
	fmt.Printf(
		"%s %s %s\n",
		colorFail.Sprint(" FAIL "),
		colorTestPath.Sprint(l.path),
		colorMessage.Sprint(msg),
	)
}

func (l *stdLog) failSQL(msg, s string, i int) {
	l.fail(msg)
	fmt.Printf(
		"\n     @ %s\n     | %s\n\n",
		colorFailPath.Sprint(pathByIndex(s, i)),
		colorText.Sprint(textByIndex(s, i)),
	)
}

func pathByIndex(s string, i int) string {
	endCount := 0
	for j := strings.LastIndex(s[:i], lineBreak); j >= 0; j-- {
		prefixEnd := j + len(fileStartPrefix)
		if s[j:prefixEnd] == fileStartPrefix {
			if endCount > 0 {
				endCount--
				continue
			}
			return s[prefixEnd : prefixEnd+strings.Index(s[prefixEnd:], lineBreak)]
		}
		if s[j:j+len(fileEndPrefix)] == fileEndPrefix {
			endCount++
		}
	}
	panic("invalid file")
}

func textByIndex(s string, i int) string {
	return s[strings.LastIndex(s[:i], lineBreak)+1:i] +
		colorFailChar.Sprint(s[i:i+1]) +
		s[i+1:i+strings.Index(s[i:], lineBreak)]
}
