package main

// factory pattern to get user
//Strategy Pattern – To switch between Human vs AI players dynamically.
import (
	"fmt"
)

// Board struct
type Board struct {
	grid [3][3]string
}

// NewBoard initializes an empty board
func NewBoard() *Board {
	return &Board{
		grid: [3][3]string{},
	}
}

// Display prints the board
func (b *Board) Display() {
	for _, row := range b.grid {
		for _, cell := range row {
			if cell == "" {
				fmt.Print("_ ")
			} else {
				fmt.Print(cell + " ")
			}
		}
		fmt.Println()
	}
}

// MakeMove updates the board
func (b *Board) MakeMove(x, y int, mark string) bool {
	if x < 0 || x >= 3 || y < 0 || y >= 3 || b.grid[x][y] != "" {
		return false
	}
	b.grid[x][y] = mark
	return true
}

// CheckWinner returns the winner, if any
func (b *Board) CheckWinner() string {
	lines := [][][2]int{
		{{0, 0}, {0, 1}, {0, 2}}, {{1, 0}, {1, 1}, {1, 2}}, {{2, 0}, {2, 1}, {2, 2}}, // Rows
		{{0, 0}, {1, 0}, {2, 0}}, {{0, 1}, {1, 1}, {2, 1}}, {{0, 2}, {1, 2}, {2, 2}}, // Columns
		{{0, 0}, {1, 1}, {2, 2}}, {{0, 2}, {1, 1}, {2, 0}}, // Diagonals
	}
	for _, line := range lines {
		if b.grid[line[0][0]][line[0][1]] != "" &&
			b.grid[line[0][0]][line[0][1]] == b.grid[line[1][0]][line[1][1]] &&
			b.grid[line[1][0]][line[1][1]] == b.grid[line[2][0]][line[2][1]] {
			return b.grid[line[0][0]][line[0][1]]
		}
	}
	return ""
}

// Player interface
type Player interface {
	GetMove(*Board) (int, int)
	GetSymbol() string
}

// HumanPlayer struct
type HumanPlayer struct {
	symbol string
}

// GetMove prompts the user for input
func (p *HumanPlayer) GetMove(b *Board) (int, int) {
	var x, y int
	fmt.Println("Enter row and column (0-2):")
	fmt.Scan(&x, &y)
	return x, y
}

// GetSymbol returns the player's symbol
func (p *HumanPlayer) GetSymbol() string {
	return p.symbol
}

// AIPlayer struct (Simple AI for demonstration)
type AIPlayer struct {
	symbol string
}

// GetMove returns the first available move
func (p *AIPlayer) GetMove(b *Board) (int, int) {
	for i := range b.grid {
		for j := range b.grid[i] {
			if b.grid[i][j] == "" {
				return i, j
			}
		}
	}
	return -1, -1
}

// GetSymbol returns the AI's symbol
func (p *AIPlayer) GetSymbol() string {
	return p.symbol
}

// PlayerFactory to create players dynamically
func PlayerFactory(playerType, symbol string) Player {
	if playerType == "human" {
		return &HumanPlayer{symbol: symbol}
	} else if playerType == "ai" {
		return &AIPlayer{symbol: symbol}
	}
	return nil
}

// Game struct
type Game struct {
	board   *Board
	player1 Player
	player2 Player
}

// NewGame initializes the game
func NewGame(p1, p2 Player) *Game {
	return &Game{
		board:   NewBoard(),
		player1: p1,
		player2: p2,
	}
}

// Play runs the game loop
func (g *Game) Play() {
	currentPlayer := g.player1
	for {
		g.board.Display()
		x, y := currentPlayer.GetMove(g.board)

		if !g.board.MakeMove(x, y, currentPlayer.GetSymbol()) {
			fmt.Println("Invalid move, try again.")
			continue
		}

		winner := g.board.CheckWinner()
		if winner != "" {
			g.board.Display()
			fmt.Printf("Player '%s' wins!\n", winner)
			break
		}

		// Switch player
		if currentPlayer == g.player1 {
			currentPlayer = g.player2
		} else {
			currentPlayer = g.player1
		}
	}
}

func main() {
	// Creating players
	player1 := PlayerFactory("human", "X")
	player2 := PlayerFactory("ai", "O")

	// Start the game
	game := NewGame(player1, player2)
	game.Play()
}
