package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/1saswata/kanban-tui-go/internal/kanban"
	"github.com/1saswata/kanban-tui-go/internal/ui"
)

func main() {
	err := os.MkdirAll("./db", 0755)
	if err != nil {
		log.Fatal(err)
	}
	s, err := kanban.NewSQLiteStore("./db/test.db")
	if err != nil {
		log.Fatal(err)
	}
	board := ui.InitBoard(s)
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
