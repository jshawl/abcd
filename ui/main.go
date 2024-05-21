package ui

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	command Command
	diff    Diff
	help    Help
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.diff.Init(), m.command.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
		m.help, _ = m.help.Update(msg)
	}

	m.command, cmd = m.command.Update(msg)
	cmds = append(cmds, cmd)

	m.diff, cmd = m.diff.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.command.isOpen {
		return m.command.View()
	}
	if m.help.isOpen {
		return m.help.View()
	}

	return fmt.Sprintf("%s %s", m.diff.View(), m.help.View())
}

func Render() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		f.Truncate(0)
		f.Seek(0, 0)
		log.Println("program starting...")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	staged := flag.Bool("staged", false, "diff --staged ?")
	flag.Parse()

	args := flag.Args()

	p := tea.NewProgram(
		model{
			help:    NewHelp(),
			diff:    NewDiff(*staged, args),
			command: NewCommand(args),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
