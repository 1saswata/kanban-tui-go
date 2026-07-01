package ui

import (
	"charm.land/bubbles/v2/list"
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
