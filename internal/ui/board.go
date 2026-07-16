package ui

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

type Board struct {
	TaskStore kanban.TaskStore
	Columns   []Column
	Focused   int8
	err       error
}

type errMsg error

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
	return fetchTasks(b.TaskStore)
}

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		b.err = error(msg)
		return b, nil
	case tasksLoadedMsg:
		tasks := make(map[kanban.Status][]list.Item)
		for _, task := range msg {
			tasks[task.Status] = append(tasks[task.Status], taskItem{task: task})
		}
		msgTODO := b.Columns[0].list.SetItems(tasks[kanban.StatusTodo])
		msgDOING := b.Columns[1].list.SetItems(tasks[kanban.StatusDoing])
		msgDONE := b.Columns[2].list.SetItems(tasks[kanban.StatusDone])
		return b, tea.Batch(msgTODO, msgDOING, msgDONE)
	case tasksUpdatedMsg:
		return b, fetchTasks(b.TaskStore)
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
		case "enter":
			item := b.Columns[b.Focused].list.SelectedItem()
			if item == nil {
				return b, nil
			}
			task, ok := item.(taskItem)
			if !ok {
				return b, func() tea.Msg {
					return errMsg(fmt.Errorf("error getting current task"))
				}
			}
			return b, moveTask(b.TaskStore, task.task.ID, task.task.Status)
		}
	}
	var cmd tea.Cmd
	b.Columns[b.Focused].list, cmd = b.Columns[b.Focused].list.Update(msg)
	return b, cmd
}

func (b *Board) View() tea.View {
	styleError := lipgloss.NewStyle().Foreground(lipgloss.Red).Bold(true).
		Padding(0, 1)
	styleFocused := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	styleUnfocused := lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	if b.err != nil {
		return tea.NewView(styleError.Render(fmt.Sprintf(
			"Fatal Error: %v\nPress q to quit", b.err)))
	}
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
