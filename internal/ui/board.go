package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

type Board struct {
	TaskStore     kanban.TaskStore
	Columns       []Column
	Focused       int8
	err           error
	input         textinput.Model
	isTyping      bool
	descInput     textarea.Model
	isEditingDesc bool
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
	b.input = textinput.New()
	b.input.Placeholder = "New Task Title"
	b.input.SetWidth(20)
	b.isTyping = false
	b.descInput = textarea.New()
	b.descInput.SetHeight(20)
	b.descInput.SetWidth(10)
	b.isEditingDesc = false
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
		cmdTODO := b.Columns[0].list.SetItems(tasks[kanban.StatusTodo])
		cmdDOING := b.Columns[1].list.SetItems(tasks[kanban.StatusDoing])
		cmdDONE := b.Columns[2].list.SetItems(tasks[kanban.StatusDone])
		return b, tea.Batch(cmdTODO, cmdDOING, cmdDONE)
	case tasksUpdatedMsg:
		return b, fetchTasks(b.TaskStore)
	case tea.KeyMsg:
		if !b.isTyping {
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
				task, ok := b.getSelectedTask()
				if !ok {
					return b, func() tea.Msg {
						return errMsg(fmt.Errorf("error getting current task"))
					}
				}
				if task.Status == kanban.StatusDone {
					return b, nil
				}
				nextStatus := task.Status
				switch task.Status {
				case kanban.StatusTodo:
					nextStatus = kanban.StatusDoing
				case kanban.StatusDoing:
					nextStatus = kanban.StatusDone
				}
				return b, setTaskStatus(b.TaskStore, task.ID, nextStatus)
			case "backspace":
				task, ok := b.getSelectedTask()
				if !ok {
					return b, func() tea.Msg {
						return errMsg(fmt.Errorf("error getting current task"))
					}
				}
				if task.Status == kanban.StatusTodo {
					return b, nil
				}
				nextStatus := task.Status
				switch task.Status {
				case kanban.StatusDoing:
					nextStatus = kanban.StatusTodo
				case kanban.StatusDone:
					nextStatus = kanban.StatusDoing
				}
				return b, setTaskStatus(b.TaskStore, task.ID, nextStatus)
			case "n":
				b.isTyping = true
				return b, b.input.Focus()
			case "x", "delete":
				task, ok := b.getSelectedTask()
				if !ok {
					return b, func() tea.Msg {
						return errMsg(fmt.Errorf("error getting current task"))
					}
				}
				return b, deleteTask(b.TaskStore, task.ID)
			}
		} else {
			switch msg.String() {
			case "enter":
				val := b.input.Value()
				b.isTyping = false
				b.input.Reset()
				if strings.TrimSpace(val) == "" {
					return b, nil
				}
				return b, createTask(b.TaskStore, val)
			case "esc":
				b.input.Reset()
				b.isTyping = false
				return b, nil
			default:
				var cmd tea.Cmd
				b.input, cmd = b.input.Update(msg)
				return b, cmd
			}
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
	list = lipgloss.JoinHorizontal(lipgloss.Top, list, b.renderDetailView())
	if b.isTyping {
		input := b.input.View()
		list = lipgloss.JoinVertical(lipgloss.Left, list, input)
	}
	v := tea.NewView(lipgloss.NewStyle().Margin(1, 0, 2, 4).Render(list))
	v.AltScreen = true
	return v
}

func (b *Board) getSelectedTask() (kanban.Task, bool) {
	var ti taskItem
	item := b.Columns[b.Focused].list.SelectedItem()
	if item == nil {
		return ti.task, false
	}
	ti, ok := item.(taskItem)
	return ti.task, ok
}

func (b *Board) renderDetailView() string {
	boxStyle := lipgloss.NewStyle().
		Width(30).Height(10).Border(lipgloss.NormalBorder())
	task, ok := b.getSelectedTask()
	if !ok {
		return boxStyle.Render("No task selected")
	}
	return boxStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left, "Title: "+task.Title, "Status: "+string(task.Status),
		"Description: "+task.Description))
}
