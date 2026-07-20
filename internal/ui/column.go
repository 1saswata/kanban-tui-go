package ui

import (
	"context"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/1saswata/kanban-tui-go/internal/kanban"
	"github.com/google/uuid"
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

type tasksUpdatedMsg struct{}

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
	return t.task.Title
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

func moveTask(store kanban.TaskStore, id uuid.UUID, status kanban.Status) tea.Cmd {
	return func() tea.Msg {
		if status == kanban.StatusDone {
			return nil
		}
		nextStatus := status
		switch status {
		case kanban.StatusTodo:
			nextStatus = kanban.StatusDoing
		case kanban.StatusDoing:
			nextStatus = kanban.StatusDone
		}
		err := store.UpdateTaskStatus(context.Background(), id, nextStatus)
		if err != nil {
			return errMsg(err)
		}
		return tasksUpdatedMsg{}
	}
}

func createTask(store kanban.TaskStore, title string) tea.Cmd {
	return func() tea.Msg {
		task := kanban.Task{
			ID:        uuid.New(),
			Title:     title,
			Status:    kanban.StatusTodo,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.CreateTask(context.Background(), task)
		if err != nil {
			return errMsg(err)
		}
		return tasksUpdatedMsg{}
	}
}
