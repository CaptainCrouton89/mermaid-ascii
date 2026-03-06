package cmd

import (
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var htmlBreakPattern = regexp.MustCompile(`(?i)<br\s*/?>`)

type graphLabel struct {
	lines []string
	width int
}

func newGraphLabel(raw string) graphLabel {
	normalized := htmlBreakPattern.ReplaceAllString(raw, "\n")
	normalized = strings.ReplaceAll(normalized, `\n`, "\n")

	lines := strings.Split(normalized, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}

	width := 0
	for _, line := range lines {
		width = Max(width, runewidth.StringWidth(line))
	}

	return graphLabel{
		lines: lines,
		width: width,
	}
}

func (l graphLabel) height() int {
	return len(l.lines)
}
