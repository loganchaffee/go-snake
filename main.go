package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	Left  string = "left"
	Right string = "right"
	Up    string = "up"
	Down  string = "down"

	Break string = "break"
)

type Coord struct {
	X int
	Y int
}

type Snake struct {
	Head Coord
	Tail []Coord
}

type GameState struct {
	Snake Snake
	Input string
	Food  Coord
}

func enableRawMode() {
	cmd := exec.Command("/bin/stty", "raw", "-echo")
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func disableRawMode() {
	cmd := exec.Command("/bin/stty", "-raw", "echo")
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func clearScreen() {
	// Use ANSI escape codes to clear the screen and move the cursor
	fmt.Print("\033[H\033[2J")
}

func hideCursor() {
	fmt.Print("\033[H")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func render(gameState *GameState) {
	clearScreen()

	// Create board array
	var board [20][20]string
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			if j == 0 || j == 19 || i == 0 || i == 19 {
				board[i][j] = "⬜️"
			} else {
				board[i][j] = "⬛️"
			}
		}
	}

	// Add food to board
	board[gameState.Food.Y][gameState.Food.X] = "🍕"

	// Add tail pieces to board
	for i := 0; i < len(gameState.Snake.Tail); i++ {
		tailPiece := gameState.Snake.Tail[i]

		board[tailPiece.Y][tailPiece.X] = "🟢"
	}

	// Add snake head to board
	board[gameState.Snake.Head.Y][gameState.Snake.Head.X] = "🟢"

	// Create string from board
	page := ""

	// Score
	page += "Score: " + strconv.Itoa(1+len(gameState.Snake.Tail)) + "\n\r"

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			page += board[i][j]
		}
		page += "\n\r"
	}

	fmt.Println(page)
}

func handleInput(input *string) {
	var buf [1]byte

	for {
		_, err := os.Stdin.Read(buf[:])

		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			break
		}

		switch buf[0] {
		case 'a':
			*input = Left
		case 'd':
			*input = Right
		case 'w':
			*input = Up
		case 's':
			*input = Down
		case 'q':
			*input = Break
		}

		if buf[0] == 'q' {
			break
		}
	}
}

func initGameState() GameState {
	snake := Snake{}

	tail := []Coord{
		// {3, 1},
		// {2, 1},
		// {1, 1},
	}

	snake.Tail = tail

	snake.Head = Coord{4, 1}

	input := Right

	gameState := GameState{
		Snake: snake,
		Input: input,
		Food:  Coord{4, 8},
	}

	return gameState
}

func checkIsEndOfGame(gameState *GameState) bool {
	// hit wall
	if gameState.Snake.Head.X == 19 || gameState.Snake.Head.X == 0 || gameState.Snake.Head.Y == 19 || gameState.Snake.Head.Y == 0 {
		return true
	}

	// entered escape key ("q")
	if gameState.Input == Break {
		return true
	}

	for i := 0; i < len(gameState.Snake.Tail); i++ {
		tailPiece := gameState.Snake.Tail[i]
		if tailPiece == gameState.Snake.Head {
			return true
		}
	}

	return false
}

func setFoodCoord(gameState *GameState) {
	available := make([]Coord, 1)

	for i := 1; i < 19; i++ {
		for j := 1; j < 19; j++ {
			currentCoord := Coord{i, j}

			isInHead := false
			isInTail := false

			if gameState.Snake.Head == currentCoord {
				isInHead = true
			}

			for k := 0; k < len(gameState.Snake.Tail); k++ {
				if gameState.Snake.Tail[k] == currentCoord {
					isInTail = true
				}
			}

			if !isInHead && !isInTail {
				available = append(available, currentCoord)
			}
		}
	}

	// Generate random index
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(available))

	gameState.Food = available[randomIndex]
}

func main() {
	defer func() {
		disableRawMode()
		showCursor()
	}()

	enableRawMode()

	gameState := initGameState()

	stagedInput := ""

	go handleInput(&stagedInput)

	for {
		render(&gameState)

		if checkIsEndOfGame(&gameState) {
			break
		}

		switch stagedInput {
		case Left:
			if gameState.Input != Right {
				gameState.Input = stagedInput
			}
		case Right:
			if gameState.Input != Left {
				gameState.Input = stagedInput
			}
		case Up:
			if gameState.Input != Down {
				gameState.Input = stagedInput
			}
		case Down:
			if gameState.Input != Up {
				gameState.Input = stagedInput
			}
		case Break:
			gameState.Input = stagedInput
		}

		// Did we eat the food
		if gameState.Food == gameState.Snake.Head {

			// Is there just a head
			if len(gameState.Snake.Tail) == 0 {
				lastPiece := gameState.Snake.Head

				switch gameState.Input {
				case Left:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X + 1, lastPiece.Y})
				case Right:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X - 1, lastPiece.Y})
				case Up:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X, lastPiece.Y + 1})
				case Down:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X, lastPiece.Y - 1})
				}
			} else if len(gameState.Snake.Tail) == 1 {
				lastPiece := gameState.Snake.Tail[len(gameState.Snake.Tail)-1]

				switch gameState.Input {
				case Left:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X + 1, lastPiece.Y})
				case Right:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X - 1, lastPiece.Y})
				case Up:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X, lastPiece.Y + 1})
				case Down:
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X, lastPiece.Y - 1})
				}

			} else {
				lastPiece := gameState.Snake.Tail[len(gameState.Snake.Tail)-1]
				secondToLastPiece := gameState.Snake.Tail[len(gameState.Snake.Tail)-2]

				// Append to left
				if secondToLastPiece.X > lastPiece.X {
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X - 1, lastPiece.Y})
				}

				// Append to right
				if secondToLastPiece.X < lastPiece.X {
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X + 1, lastPiece.Y})
				}

				// Append to bottom
				if secondToLastPiece.Y < lastPiece.Y {
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X, lastPiece.Y + 1})
				}

				// Append to top
				if secondToLastPiece.X > lastPiece.X {
					gameState.Snake.Tail = append(gameState.Snake.Tail, Coord{lastPiece.X, lastPiece.Y - 1})
				}
			}

			setFoodCoord(&gameState)
		}

		// Update snake tail position
		for i := len(gameState.Snake.Tail) - 1; i > -1; i-- {
			if i == 0 {
				gameState.Snake.Tail[i] = gameState.Snake.Head
			} else {
				gameState.Snake.Tail[i] = gameState.Snake.Tail[i-1]
			}
		}

		// Update snake head position
		switch gameState.Input {
		case Right:
			gameState.Snake.Head = Coord{
				X: gameState.Snake.Head.X + 1,
				Y: gameState.Snake.Head.Y,
			}
		case Left:
			gameState.Snake.Head = Coord{
				X: gameState.Snake.Head.X - 1,
				Y: gameState.Snake.Head.Y,
			}
		case Up:
			gameState.Snake.Head = Coord{
				X: gameState.Snake.Head.X,
				Y: gameState.Snake.Head.Y - 1,
			}
		case Down:
			gameState.Snake.Head = Coord{
				X: gameState.Snake.Head.X,
				Y: gameState.Snake.Head.Y + 1,
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
