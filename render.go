package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	parser "github.com/jshawl/diffrn/parser"
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	diff     parser.Diff
}

type refreshMsg struct {
	rawDiff string
}

func (m model) Init() tea.Cmd {
	return nil
}

func buildOutput(m model) string {
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

func (m *model) windowSizeUpdate(msg tea.WindowSizeMsg) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
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

func (m *model) refreshUpdate(msg refreshMsg) {
	m.diff, _ = parser.ParseDiff(msg.rawDiff)
	m.viewport.SetContent(buildOutput(*m))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case refreshMsg:
		m.refreshUpdate(msg)
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.windowSizeUpdate(msg)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render("diffrn")
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}

func (m model) footerView() string {
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

func render() {
	p := tea.NewProgram(
		model{},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		for {
			pause := time.Duration(1) * time.Second
			p.Send(refreshMsg{rawDiff: gitDiffRaw()})
			time.Sleep(pause)
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
