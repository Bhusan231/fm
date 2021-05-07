package ui

import (
	"io/fs"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/constants"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Files             []fs.FileInfo
	PrimaryViewport   viewport.Model
	SecondaryViewport viewport.Model
	Textinput         textinput.Model
	Spinner           spinner.Model
	Cursor            int
	ScreenWidth       int
	ScreenHeight      int
	ShowCommandBar    bool
	Ready             bool
	ActivePane        string
}

func NewModel() Model {
	cfg := config.GetConfig()

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Width = 50

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Components.Spinner))

	return Model{
		Files:             make([]fs.FileInfo, 0),
		PrimaryViewport:   viewport.Model{},
		SecondaryViewport: viewport.Model{},
		Textinput:         input,
		Spinner:           s,
		Cursor:            0,
		ScreenWidth:       0,
		ScreenHeight:      0,
		ShowCommandBar:    false,
		Ready:             false,
		ActivePane:        constants.PrimaryPane,
	}
}
