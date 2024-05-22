package ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Help struct {
	isOpen bool
	keys   map[string]string
}

func NewHelp() Help {
	return Help{
		isOpen: false,
		keys: map[string]string{
			"?":   "toggle help",
			"q":   "quit",
			"s":   "toggle --staged",
			"c":   "change git diff command arguments",
			"tab": "jump to next file",
		},
	}
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
	var content strings.Builder
	if m.isOpen {
		keys := make([]string, 0)
		for key := range m.keys {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			content.WriteString(fmt.Sprintf("%s\t%s\n", key, m.keys[key]))
		}
		return content.String()
	} else {
		return "? toggle help"
	}
}
