package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Help struct {
	isOpen bool
}

func (m Help) Init() tea.Cmd {
	return nil
}

func (m Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "?" {
			m.isOpen = !m.isOpen
		} else {
			m.isOpen = false
		}
	}
	return m, nil
}

func (m Help) View() string {
	if m.isOpen {
		return "? toggle help\nq quit\ns toggle --staged"
	} else {
		return "? toggle help"
	}
}
