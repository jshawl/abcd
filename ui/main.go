package ui

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	diff Diff
	help Help
}

func (m model) Init() tea.Cmd {
	return m.diff.Init()
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

	m.diff, cmd = m.diff.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.help.isOpen {
		return m.help.View()
	}

	return fmt.Sprintf("%s %s", m.diff.View(), m.help.View())
}

func Render() {
	if len(os.Getenv("DEBUG")) > 0 {
		file, err := tea.LogToFile("debug.log", "debug")
		_ = file.Truncate(0)
		_, _ = file.Seek(0, 0)

		log.Println("program starting...")

		if err != nil {
			fmt.Println("fatal:", err) //nolint:forbidigo

			defer func() {
				os.Exit(1)
			}()
		}

		defer file.Close()
	}

	staged := flag.Bool("staged", false, "diff --staged ?")
	flag.Parse()

	program := tea.NewProgram(
		model{
			help: NewHelp(),
			diff: NewDiff(*staged, flag.Args()),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Println("could not run program:", err) //nolint:forbidigo

		defer func() {
			os.Exit(1)
		}()
	}
}
