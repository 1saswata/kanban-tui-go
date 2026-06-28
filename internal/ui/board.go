package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/1saswata/kanban-tui-go/internal/kanban"
)

type Board struct {
	TaskStore *kanban.TaskStore
}

func (b *Board) Init() tea.Cmd {

}

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

}

func (b *Board) View() string {

}
