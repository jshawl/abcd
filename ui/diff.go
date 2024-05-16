package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	parser "github.com/jshawl/diffrn/parser"
)

type Diff struct {
	content  string
	ready    bool
	viewport viewport.Model
	diff     parser.Diff
}

type refreshMsg struct {
	rawDiff string
}

func (m Diff) Init() tea.Cmd {
	return nil
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

func (m *Diff) windowSizeUpdate(msg tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height("\n")
	footerHeight := lipgloss.Height("\n")
	verticalMarginHeight := headerHeight + footerHeight

	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.SetContent(m.content)
		m.ready = true
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
}

func (m *Diff) refreshUpdate(msg refreshMsg) {
	m.diff, _ = parser.ParseDiff(msg.rawDiff)
	m.viewport.SetContent(buildOutput(*m))
}

func (m Diff) Update(msg tea.Msg) (Diff, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case refreshMsg:
		m.refreshUpdate(msg)

	case tea.WindowSizeMsg:
		m.windowSizeUpdate(msg)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Diff) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m Diff) headerView() string {
	title := titleStyle.Render("diffrn")
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}

func (m Diff) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func gitDiffRaw() string {
	cmd := exec.Command("git", "diff")
	stdout, _ := cmd.Output()
	return string(stdout)
}
