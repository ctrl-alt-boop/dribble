package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
)

var (
	movingBlockFrames = []string{
		"██        ", " ██       ", "  ██      ", "   ██     ", "    ██    ", "     ██   ", "      ██  ", "       ██ ",
		"        ██", "       ██ ", "      ██  ", "     ██   ", "    ██    ", "   ██     ", "  ██      ", " ██       "}
	MovingBlock = spinner.Spinner{
		Frames: movingBlockFrames,
		FPS:    time.Second / 18, //nolint:mnd
	}

	GrowingBlock = spinner.Spinner{
		Frames: []string{"█", "███", "█████", "███████", "█████████", "███████", "█████", "███", "█"},
		FPS:    time.Second / 10, //nolint:mnd
	}
)
