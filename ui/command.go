package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Command struct {
	isOpen bool
	form   *huh.Form
	cmd    string
}

type commandMsg string

func (m Command) Init() tea.Cmd {
	return m.form.Init()
}

func NewCommand(args []string) Command {
	cur := strings.Join(args, " ")
	m := Command{}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("cmd").
				Title("Set command:").
				Prompt("git diff ").
				Value(&cur),
		),
	)
	return m
}

func (m Command) Update(msg tea.Msg) (Command, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if m.isOpen {
			if k == "enter" {
				m.isOpen = false
			}
		} else {
			if k == "c" {
				m.isOpen = !m.isOpen
			}
			return m, nil
		}
	}
	form, cmd := m.form.Update(msg)
	cmds = append(cmds, cmd)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		gitDiffArgs := f.GetString("cmd")
		if m.form.State == huh.StateCompleted {
			m = NewCommand([]string{gitDiffArgs})
			nextCommand := func() tea.Msg { return commandMsg(gitDiffArgs) }
			cmds = append(cmds, m.Init(), nextCommand)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Command) View() string {
	return m.form.View()
}
