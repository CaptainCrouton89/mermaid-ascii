package cmd

import (
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

func (l graphLabel) withWrapWidth(width int) graphLabel {
	if width <= 0 {
		return l
	}
	var wrapped []string
	for _, line := range l.lines {
		if runewidth.StringWidth(line) <= width {
			wrapped = append(wrapped, line)
			continue
		}
		// Simple word-wrap
		words := strings.Fields(line)
		current := ""
		currentWidth := 0
		for _, word := range words {
			wordWidth := runewidth.StringWidth(word)
			if current == "" {
				current = word
				currentWidth = wordWidth
				continue
			}
			if currentWidth+1+wordWidth <= width {
				current += " " + word
				currentWidth += 1 + wordWidth
			} else {
				wrapped = append(wrapped, current)
				current = word
				currentWidth = wordWidth
			}
		}
		if current != "" {
			wrapped = append(wrapped, current)
		}
	}
	if len(wrapped) == 0 {
		wrapped = []string{""}
	}
	maxW := 0
	for _, line := range wrapped {
		if w := runewidth.StringWidth(line); w > maxW {
			maxW = w
		}
	}
	return graphLabel{lines: wrapped, width: maxW}
}

var htmlBreakPattern = regexp.MustCompile(`(?i)<br\s*/?>`)

const graphLabelLineGap = 1

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

func (l graphLabel) contentHeight() int {
	if len(l.lines) == 0 {
		return 0
	}
	return len(l.lines) + (len(l.lines)-1)*graphLabelLineGap
}
