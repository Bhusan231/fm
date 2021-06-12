package ui

import (
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case directoryMsg:
		m.showCommandBar = false
		m.dirTree.SetContent(msg)
		m.dirTree.GotoTop()
		m.textInput.Blur()
		m.textInput.Reset()
		selectedFile, status, fileTotals, logo := m.getStatusBarContent()
		m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
		m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))

		return m, cmd

	case fileContentMsg:
		m.activeMarkdownSource = string(msg.markdownContent)
		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg.fileContent)))

		return m, cmd

	case markdownMsg:
		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg)))

		return m, cmd

	case tea.WindowSizeMsg:
		cfg := config.GetConfig()

		if !m.ready {
			m.primaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				true,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))

			m.secondaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				false,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)

			selectedFile, status, fileTotals, logo := m.getStatusBarContent()
			m.statusBar = statusbar.NewModel(
				msg.Width,
				selectedFile,
				status,
				fileTotals,
				logo,
				statusbar.Color{
					Background: cfg.Colors.StatusBar.SelectedFile.Background,
					Foreground: cfg.Colors.StatusBar.SelectedFile.Foreground,
				},
				statusbar.Color{
					Background: cfg.Colors.StatusBar.Bar.Background,
					Foreground: cfg.Colors.StatusBar.Bar.Foreground,
				},
				statusbar.Color{
					Background: cfg.Colors.StatusBar.TotalFiles.Background,
					Foreground: cfg.Colors.StatusBar.TotalFiles.Foreground,
				},
				statusbar.Color{
					Background: cfg.Colors.StatusBar.Logo.Background,
					Foreground: cfg.Colors.StatusBar.Logo.Foreground,
				},
			)

			m.statusBar.SetContent(selectedFile, status, fileTotals, logo)

			m.ready = true
		} else {
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetSize(msg.Width)
		}

		if m.activeMarkdownSource != "" {
			return m, renderMarkdownContent(m.secondaryPane.Width, m.activeMarkdownSource)
		}

		return m, cmd

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoUp()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.LineUp(3)
				}
			}

			return m, cmd

		case tea.MouseWheelDown:
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoDown()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.LineDown(3)
				}
			}

			return m, cmd
		}

	case tea.KeyMsg:
		if msg.String() == "g" && m.previousKey.String() == "g" {
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.previousKey = tea.KeyMsg{}
					m.dirTree.GotoTop()
					m.primaryPane.GotoTop()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.GotoTop()
				}
			}

			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if !m.showCommandBar {
				return m, tea.Quit
			}

		case "left", "h":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					previousPath, err := os.Getwd()

					if err != nil {
						log.Fatal("error getting working directory")
					}

					m.previousDirectory = previousPath

					return m, updateDirectoryListing(constants.PreviousDirectory, m.dirTree.ShowHidden)
				}
			}

		case "down", "j":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoDown()
					m.scrollPrimaryPane()
					selectedFile, status, fileTotals, logo := m.getStatusBarContent()
					m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.LineDown(1)
				}
			}

		case "up", "k":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoUp()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
					selectedFile, status, fileTotals, logo := m.getStatusBarContent()
					m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
				} else {
					m.secondaryPane.LineUp(1)
				}
			}

		case "G":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GotoBottom()
					m.primaryPane.GotoBottom()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.GotoBottom()
				}
			}

		case "right", "l":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					if m.dirTree.GetSelectedFile().IsDir() && !m.textInput.Focused() {
						return m, updateDirectoryListing(m.dirTree.GetSelectedFile().Name(), m.dirTree.ShowHidden)
					} else {
						m.secondaryPane.GotoTop()

						return m, m.readFileContent(m.dirTree.GetSelectedFile())
					}
				}
			}

		case "enter":
			command, value := utils.ParseCommand(m.textInput.Value())

			if command == "" {
				return m, nil
			}

			switch command {
			case "mkdir":
				return m, createDir(value, m.dirTree.ShowHidden)

			case "touch":
				return m, createFile(value, m.dirTree.ShowHidden)

			case "mv", "rename":
				return m, renameFileOrDir(m.dirTree.GetSelectedFile().Name(), value, m.dirTree.ShowHidden)

			case "cp":
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, moveDir(m.dirTree.GetSelectedFile().Name(), value, m.dirTree.ShowHidden)
				} else {
					return m, moveFile(m.dirTree.GetSelectedFile().Name(), value, m.dirTree.ShowHidden)
				}

			case "rm", "delete":
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, deleteDir(m.dirTree.GetSelectedFile().Name(), m.dirTree.ShowHidden)
				} else {
					return m, deleteFile(m.dirTree.GetSelectedFile().Name(), m.dirTree.ShowHidden)
				}

			default:
				return m, nil
			}

		case ":":
			m.showCommandBar = true
			m.textInput.Placeholder = "enter command"
			m.textInput.Focus()

			return m, cmd

		case "~":
			if !m.showCommandBar {
				return m, updateDirectoryListing(utils.GetHomeDirectory(), m.dirTree.ShowHidden)
			}

		case "-":
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, updateDirectoryListing(m.previousDirectory, m.dirTree.ShowHidden)
			}

		case ".":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.ToggleHidden()

				return m, updateDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
			}

		case "tab":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.primaryPane.IsActive = false
					m.secondaryPane.IsActive = true
				} else {
					m.primaryPane.IsActive = true
					m.secondaryPane.IsActive = false
				}
			}

		case "esc":
			m.showCommandBar = false
			m.textInput.Blur()
			m.textInput.Reset()
			m.secondaryPane.GotoTop()
			m.primaryPane.IsActive = true
			m.secondaryPane.IsActive = false
			selectedFile, status, fileTotals, logo := m.getStatusBarContent()
			m.statusBar.SetContent(selectedFile, status, fileTotals, logo)

			return m, renderMarkdownContent(m.secondaryPane.Width, constants.HelpText)
		}

		m.previousKey = msg
	}

	if m.showCommandBar {
		selectedFile, status, fileTotals, logo := m.getStatusBarContent()
		m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
