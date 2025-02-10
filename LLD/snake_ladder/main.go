package main

import (
	"fmt"
	"math/rand"
	"time"
)

// --- Strategy Pattern for Dice Rolling ---
type Dice interface {
	Roll() int
}

type NormalDice struct{}

func (d *NormalDice) Roll() int {
	return rand.Intn(6) + 1
}

type CrookedDice struct{}

func (d *CrookedDice) Roll() int {
	return 2 * (rand.Intn(3) + 1)
}

// --- Factory Pattern for Board Components ---
type BoardComponent interface {
	AffectPosition(int) int
}

type Snake struct {
	start, end int
}

func (s *Snake) AffectPosition(pos int) int {
	if pos == s.start {
		return s.end
	}
	return pos
}

type Ladder struct {
	start, end int
}

func (l *Ladder) AffectPosition(pos int) int {
	if pos == l.start {
		return l.end
	}
	return pos
}

// BoardComponent Factory
func NewBoardComponent(componentType string, start, end int) BoardComponent {
	if componentType == "snake" {
		return &Snake{start, end}
	} else if componentType == "ladder" {
		return &Ladder{start, end}
	}
	return nil
}

// --- Builder Pattern for Game Board ---
type Board struct {
	size       int
	components []BoardComponent
}
type IBoardBuilder interface {
	AddComponent(component BoardComponent) IBoardBuilder
	Build() *Board
}

type BoardBuilder struct {
	board Board
}

func NewBoardBuilder(size int) *BoardBuilder {
	return &BoardBuilder{board: Board{size: size}}
}

func (bb *BoardBuilder) AddComponent(component BoardComponent) IBoardBuilder {
	bb.board.components = append(bb.board.components, component)
	return bb
}

func (bb *BoardBuilder) Build() *Board {
	return &bb.board
}

// --- User Class for Making Moves ---
type User struct {
	name     string
	position int
	dice     Dice
}

func (u *User) Move(board *Board) {
	steps := u.dice.Roll()
	fmt.Printf("%s rolled a %d\n", u.name, steps)
	newPos := u.position + steps
	if newPos > board.size {
		return
	}

	for _, component := range board.components {
		newPos = component.AffectPosition(newPos)
	}

	u.position = newPos
	fmt.Printf("%s moved to %d\n", u.name, u.position)
}

// --- Game Logic ---
type Game struct {
	users []*User
	board *Board
}

func NewGame(board *Board, users []string) *Game {
	var u []*User
	for _, name := range users {
		u = append(u, &User{name: name, position: 0, dice: &NormalDice{}})
	}
	return &Game{users: u, board: board}
}

func (g *Game) Play() {
	rand.Seed(time.Now().UnixNano())
	winner := false

	for !winner {
		for _, user := range g.users {
			user.Move(g.board)

			if user.position == g.board.size {
				fmt.Printf("%s wins the game!\n", user.name)
				winner = true
				break
			}
		}
	}
}

// --- Main Execution ---
func main() {
	board := NewBoardBuilder(100).
		AddComponent(NewBoardComponent("snake", 14, 7)).
		AddComponent(NewBoardComponent("ladder", 3, 22)).
		Build()

	game := NewGame(board, []string{"Alice", "Bob"})
	game.Play()
}
