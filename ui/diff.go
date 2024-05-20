package ui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	parser "github.com/jshawl/abcd/parser"
)

type Diff struct {
	ready       bool
	viewport    viewport.Model
	staged      bool
	files       []File
	currentFile int
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

func (m Diff) ViewEmpty() string {
	if m.staged {
		return "No changes staged..."
	} else {
		return "No diff to show! Working directory is clean."
	}
}

func (m Diff) lines() string {
	if len(m.files) == 0 {
		return m.ViewEmpty()
	} else {
		var content strings.Builder
		for _, file := range m.files {
			content.WriteString(file.View(m.viewport.Width))
		}
		return content.String()
	}
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
	case tea.KeyMsg:
		k := msg.String()
		if k == "s" {
			m.staged = !m.staged
			cmd = m.Tick(true)
			cmds = append(cmds, cmd)
		}
		if k == "tab" {
			heights := []int{0}
			for i := 0; i < len(m.files); i++ {
				height := lipgloss.Height(m.files[i].View(m.viewport.Width)) - 1
				total := height + heights[len(heights)-1]
				heights = append(heights, total)
			}

			m.currentFile += 1
			if m.currentFile == len(m.files) {
				m.currentFile = 0
			}
			m.viewport.SetYOffset(heights[m.currentFile])
		}
	case TickMsg:
		diff, _ := parser.ParseDiff(m.gitDiffRaw())
		m.files = []File{}
		for _, file := range diff.Files {
			m.files = append(m.files, NewFile(file))
		}
		m.viewport.SetContent(m.lines())
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
	help := "? toggle help "
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	var cmd string
	if m.staged {
		cmd = "git diff --staged"
	} else {
		cmd = "git diff"
	}
	space := m.viewport.Width - lipgloss.Width(info) - lipgloss.Width(help) - lipgloss.Width(cmd)
	line := strings.Repeat(" ", max(0, space))
	return lipgloss.JoinHorizontal(lipgloss.Center, cmd, line, help, info)
}

func (m Diff) gitDiffRaw() string {
	var cmd *exec.Cmd
	if m.staged {
		cmd = exec.Command("git", "diff", "--staged")
	} else {
		cmd = exec.Command("git", "diff")
	}
	stdout, _ := cmd.Output()
	return string(stdout)
}
