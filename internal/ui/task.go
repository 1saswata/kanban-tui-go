package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

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
