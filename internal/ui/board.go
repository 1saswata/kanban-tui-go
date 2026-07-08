package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

type Board struct {
	TaskStore kanban.TaskStore
	Columns   []Column
	Focused   int8
}

func InitBoard(ts kanban.TaskStore) *Board {
	b := &Board{TaskStore: ts}
	b.Columns = []Column{
		NewColumn(kanban.StatusTodo),
		NewColumn(kanban.StatusDoing),
		NewColumn(kanban.StatusDone),
	}
	b.Focused = 0
	return b
}

func (b *Board) Init() tea.Cmd {
	return nil
}

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			b.Focused = (b.Focused - 1 + 3) % 3
			return b, nil
		case "l", "right":
			b.Focused = (b.Focused + 1 + 3) % 3
			return b, nil
		case "ctrl+c", "q", "esc":
			return b, tea.Quit
		}
	}
	var cmd tea.Cmd
	b.Columns[b.Focused].list, cmd = b.Columns[b.Focused].list.Update(msg)
	return b, cmd
}

func (b *Board) View() tea.View {
	styleFocused := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	styleUnfocused := lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	var lists [3]string
	for i, c := range b.Columns {
		if i == int(b.Focused) {
			lists[i] = styleFocused.Render("> " + c.list.View())
		} else {
			lists[i] = styleUnfocused.Render(c.list.View())
		}
	}
	list := lipgloss.JoinHorizontal(lipgloss.Top, lists[0], lists[1], lists[2])
	v := tea.NewView(lipgloss.NewStyle().Margin(1, 0, 2, 4).Render(list))
	v.AltScreen = true
	return v
}
