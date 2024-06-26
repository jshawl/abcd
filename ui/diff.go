package ui

import (
	"fmt"
	"os/exec"
	"slices"
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
	args        []string
}

type TickMsg struct{}

func (m Diff) command() []string {
	cmd := []string{"git", "diff"}
	if m.staged {
		cmd = append(cmd, "--staged")
	}
	cmd = append(cmd, m.args...)
	return cmd
}

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

func NewDiff(staged bool, args []string) Diff {
	return Diff{staged: staged, args: args}
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

func (m Diff) FileHeights() []int {
	heights := []int{0}
	for i := range m.files {
		height := m.files[i].Height(m.viewport.Width) - 1
		total := height + heights[len(heights)-1]
		heights = append(heights, total)
	}
	return heights
}

func (m Diff) getFileIndexInViewport() int {
	index := slices.IndexFunc(m.FileHeights(), func(i int) bool {
		offset := m.viewport.YOffset - 1
		return offset-i <= -2
	}) - 1

	if index >= 0 && len(m.files) > 0 {
		return index
	}
	return m.currentFile
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
		if k == "shift+tab" {
			// decrement the current file only if the first line of the
			// file is visible in the viewport. If the viewport has scrolled
			// past the top of the file, shift+tab should scroll to the top
			// of the file.
			if slices.Contains(m.FileHeights(), m.viewport.YOffset) {
				m.currentFile -= 1
				if m.currentFile < 0 || m.viewport.AtTop() {
					m.currentFile = len(m.files) - 1
				}
			}
			m.viewport.SetYOffset(m.FileHeights()[m.currentFile])
		}
		if k == "tab" {
			m.currentFile += 1
			if m.currentFile == len(m.files) || m.viewport.AtBottom() {
				m.currentFile = 0
			}
			m.viewport.SetYOffset(m.FileHeights()[m.currentFile])
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

	m.currentFile = m.getFileIndexInViewport()

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
	cmd := strings.Join(m.command(), " ")
	space := m.viewport.Width - lipgloss.Width(info) - lipgloss.Width(help) - lipgloss.Width(cmd)
	line := strings.Repeat(" ", max(0, space))
	return lipgloss.JoinHorizontal(lipgloss.Center, cmd, line, help, info)
}

func (m Diff) gitDiffRaw() string {
	var cmd *exec.Cmd
	cmds := m.command()
	cmd = exec.Command(cmds[0], cmds[1:]...)
	stdout, _ := cmd.Output()
	return string(stdout)
}
