package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

type Board struct {
	TaskStore kanban.TaskStore
}

func InitBoard(ts kanban.TaskStore) Board {
	return Board{TaskStore: ts}
}

func (b Board) Init() tea.Cmd {
	return nil
}

func (b Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return b, tea.Quit
		}
	}
	return b, nil
}

func (b Board) View() tea.View {
	s := "Kanban Board Initializing...\nPress 'q' to quit."
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
