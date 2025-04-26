package main

import (
	"fmt"
	"strings"
)

// --------- Cursor Struct ---------
type Cursor struct {
	line, col int
}

// --------- Editor Struct ---------
type Editor struct {
	lines  []string
	cursor Cursor
}

func NewEditor() *Editor {
	return &Editor{
		lines: []string{""},
		cursor: Cursor{
			line: 0,
			col:  0,
		},
	}
}

func (e *Editor) Print() {
	for i, line := range e.lines {
		cursorIndicator := ""
		if i == e.cursor.line {
			cursorIndicator = strings.Repeat(" ", e.cursor.col) + "^"
		}
		fmt.Println(line)
		if cursorIndicator != "" {
			fmt.Println(cursorIndicator)
		}
	}
	fmt.Println()
}

// --------- Command Interface ---------
type Command interface {
	Execute(e *Editor)
}

// --------- Append Command ---------
type AppendCommand struct {
	text string
}

func (a *AppendCommand) Execute(e *Editor) {
	line := e.lines[e.cursor.line]
	before := line[:e.cursor.col]
	after := line[e.cursor.col:]
	newLine := before + a.text + after
	e.lines[e.cursor.line] = newLine
	e.cursor.col += len(a.text)
}

// --------- Replace Command ---------
type ReplaceCommand struct {
	text string
}

func (r *ReplaceCommand) Execute(e *Editor) {
	e.lines[e.cursor.line] = r.text
	e.cursor.col = len(r.text)
}

// --------- Arrow Command ---------
type ArrowCommand struct {
	direction string
}

func (a *ArrowCommand) Execute(e *Editor) {
	switch a.direction {
	case "left":
		if e.cursor.col > 0 {
			e.cursor.col--
		}
	case "right":
		if e.cursor.col < len(e.lines[e.cursor.line]) {
			e.cursor.col++
		}
	case "up":
		if e.cursor.line > 0 {
			e.cursor.line--
			if e.cursor.col > len(e.lines[e.cursor.line]) {
				e.cursor.col = len(e.lines[e.cursor.line])
			}
		}
	case "down":
		if e.cursor.line < len(e.lines)-1 {
			e.cursor.line++
			if e.cursor.col > len(e.lines[e.cursor.line]) {
				e.cursor.col = len(e.lines[e.cursor.line])
			}
		}
	}
}

// --------- Page Command ---------
type PageCommand struct {
	up bool
}

func (p *PageCommand) Execute(e *Editor) {
	pageSize := 5
	if p.up {
		e.cursor.line -= pageSize
		if e.cursor.line < 0 {
			e.cursor.line = 0
		}
	} else {
		e.cursor.line += pageSize
		if e.cursor.line >= len(e.lines) {
			for len(e.lines) <= e.cursor.line {
				e.lines = append(e.lines, "")
			}
		}
	}
	if e.cursor.col > len(e.lines[e.cursor.line]) {
		e.cursor.col = len(e.lines[e.cursor.line])
	}
}

// --------- Command Executor (Strategy Context) ---------
type CommandExecutor struct{}

func (c *CommandExecutor) ExecuteCommand(command Command, editor *Editor) {
	command.Execute(editor)
}

// --------- Main ---------
func main() {
	editor := NewEditor()
	executor := &CommandExecutor{}

	executor.ExecuteCommand(&AppendCommand{text: "Hello"}, editor)
	executor.ExecuteCommand(&ArrowCommand{direction: "right"}, editor)
	executor.ExecuteCommand(&ArrowCommand{direction: "right"}, editor)
	executor.ExecuteCommand(&AppendCommand{text: " World"}, editor)

	// Add new lines for testing PageDown
	for i := 0; i < 10; i++ {
		editor.lines = append(editor.lines, fmt.Sprintf("Line %d", i+1))
	}

	executor.ExecuteCommand(&PageCommand{up: false}, editor) // Page Down
	executor.ExecuteCommand(&ArrowCommand{direction: "down"}, editor)
	executor.ExecuteCommand(&ReplaceCommand{text: "Replaced Text"}, editor)

	editor.Print()
}
