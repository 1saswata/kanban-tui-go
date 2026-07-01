package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

type Board struct {
	TaskStore kanban.TaskStore
	Columns   []Column
}

func InitBoard(ts kanban.TaskStore) *Board {
	b := &Board{TaskStore: ts}
	b.Columns = []Column{
		NewColumn(kanban.StatusTodo),
		NewColumn(kanban.StatusDoing),
		NewColumn(kanban.StatusDone),
	}
	return b
}

func (b *Board) Init() tea.Cmd {
	return nil
}

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return b, tea.Quit
		}
	}
	return b, nil
}

func (b *Board) View() tea.View {
	lists := lipgloss.JoinHorizontal(lipgloss.Top, b.Columns[0].list.View(),
		b.Columns[1].list.View(), b.Columns[2].list.View())
	v := tea.NewView(lists)
	v.AltScreen = true
	return v
}
