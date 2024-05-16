package ui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	parser "github.com/jshawl/diffrn/parser"
)

type Diff struct {
	ready    bool
	viewport viewport.Model
	diff     parser.Diff
}

type TickMsg struct{}

func (m Diff) Tick(immediately ...bool) tea.Cmd {
	if len(immediately) > 0 {
		return func() tea.Msg {
			return TickMsg{}
		}
	}
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m Diff) Init() tea.Cmd {
	return m.Tick()
}

func buildOutput(m Diff) string {
	var content strings.Builder
	if len(m.diff.Files) == 0 {
		return "No diff to show! Working directory is clean."
	}

	for _, file := range m.diff.Files {
		content.WriteString(fileStyle.Width(m.viewport.Width).Render(file.Name))
		content.WriteString("\n")
		for blockI, block := range file.Blocks {
			for _, line := range block.Lines {
				if strings.HasPrefix(line, "-") {
					content.WriteString(removedStyle.Width(m.viewport.Width).Render(line))
				} else if strings.HasPrefix(line, "+") {
					content.WriteString(addedStyle.Width(m.viewport.Width).Render(line))
				} else {
					content.WriteString(line)
				}
				content.WriteString("\n")
			}
			if blockI < len(file.Blocks)-1 {
				content.WriteString(hr.Render("···"))
				content.WriteString("\n")
			}
		}
	}
	return content.String()
}

func (m *Diff) windowSizeUpdate(msg tea.WindowSizeMsg) tea.Cmd {
	footerHeight := lipgloss.Height(m.footerView())
	helpHeight := lipgloss.Height("\n")
	verticalMarginHeight := footerHeight + helpHeight - 2

	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = 0
		m.ready = true
		return m.Tick(true)
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		return nil
	}
}

func (m Diff) Update(msg tea.Msg) (Diff, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case TickMsg:
		m.diff, _ = parser.ParseDiff(gitDiffRaw())
		m.viewport.SetContent(buildOutput(m))
		return m, m.Tick()

	case tea.WindowSizeMsg:
		cmd = m.windowSizeUpdate(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Diff) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s", m.viewport.View(), m.footerView())
}

func (m Diff) footerView() string {
	return infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
}

func gitDiffRaw() string {
	cmd := exec.Command("git", "diff")
	stdout, _ := cmd.Output()
	return string(stdout)
}
