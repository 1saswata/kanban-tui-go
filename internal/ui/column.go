package ui

import (
	"context"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

const listHeight = 20
const listWidth = 30

type Column struct {
	list   list.Model
	status kanban.Status
}

func NewColumn(status kanban.Status) Column {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), listWidth, listHeight)
	l.Title = string(status)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	return Column{list: l, status: status}
}

type tasksLoadedMsg []kanban.Task

type taskItem struct {
	task kanban.Task
}

func (t taskItem) Title() string {
	return t.task.Title
}

func (t taskItem) Description() string {
	return t.task.Description
}

func (t taskItem) FilterValue() string {
	return ""
}

func fetchTasks(store kanban.TaskStore) tea.Cmd {
	return func() tea.Msg {
		t, err := store.GetTasks(context.Background())
		if err != nil {
			return errMsg(err)
		}
		return tasksLoadedMsg(t)
	}
}
