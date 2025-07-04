package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	charmansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/reflow/ansi"
)

// PlaceOverlay places fg on top of bg.
func PlaceOverlay(x, y lipgloss.Position, fg, bg string) string {
	if x < 0 || y < 0 || x > 1 || y > 1 {
		return bg
	}

	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return fg
	}

	topLeftX := int(math.Round(float64(bgWidth-fgWidth) * float64(x)))
	topLeftY := int(math.Round(float64(bgHeight-fgHeight) * float64(y)))

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < topLeftY || i >= topLeftY+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		left := charmansi.Truncate(bgLine, topLeftX, "")
		pos = ansi.PrintableRuneWidth(left)
		b.WriteString(left)

		if pos < topLeftX {
			b.WriteString(strings.Repeat(" ", topLeftX-pos))
			pos = topLeftX
		}

		fgLine := fgLines[i-topLeftY]
		b.WriteString(fgLine)
		pos += ansi.PrintableRuneWidth(fgLine)

		right := charmansi.TruncateLeft(bgLine, pos, "")
		bgWidth := ansi.PrintableRuneWidth(bgLine)
		rightWidth := ansi.PrintableRuneWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(strings.Repeat(" ", bgWidth-rightWidth-pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

func PlaceOverlaySimple(x, y lipgloss.Position, fg, bg string) string {
	if x < 0 || y < 0 || x > 1 || y > 1 {
		return bg
	}

	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return fg
	}

	topLeftX := int(math.Round(float64(bgWidth-fgWidth) * float64(x)))
	topLeftY := int(math.Round(float64(bgHeight-fgHeight) * float64(y)))

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < topLeftY || i >= topLeftY+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		overlayLine := charmansi.Truncate(bgLine, topLeftX, "")
		overlayLine += fgLines[i-topLeftY]
		overlayLine += charmansi.TruncateLeft(bgLine, ansi.PrintableRuneWidth(overlayLine), "")
		b.WriteString(overlayLine)
	}
	return b.String()
}

func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := charmansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
