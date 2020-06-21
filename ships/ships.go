package ships

import (
	"encoding/gob"
	"encoding/json"
	"math/rand"
	"net/http"

	"battleship.go/sessions"
)

const fromLeftToRight = 0
const fromTopToBottom = 1
const boardSize = 10
const freeSpaceIndex = 99

var shipSizes = [8]int{5, 4, 3, 3, 2, 2, 1, 1}

type singleShip struct {
	Size, Hits, Index int
}
type boardArray [][]int
type shipsArmy struct {
	Ships []singleShip
	Board boardArray
}
type BlackAndWhiteArmies struct {
	White, Black shipsArmy
}

func Init(writer http.ResponseWriter, request *http.Request) {
	gob.Register(BlackAndWhiteArmies{})

	session := sessions.CheckSession(writer, request)
	ships, success := initShips()
	if !success {
		msg := struct {
			Command string
			Message string
		}{
			Command: "display",
			Message: "Could not place all ships on the board in 100 rounds",
		}
		json.NewEncoder(writer).Encode(msg)
		return
	}

	session.Values["ships"] = ships
	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func initBoard() boardArray {
	board := make(boardArray, boardSize)
	for i := range board {
		board[i] = make([]int, boardSize)
	}
	return board
}

func initShips() (BlackAndWhiteArmies, bool) {
	white := shipsArmy{}
	black := shipsArmy{}
	if !white.setRandomPositions() {
		return BlackAndWhiteArmies{}, false
	}
	if !black.setRandomPositions() {
		return BlackAndWhiteArmies{}, false
	}
	return BlackAndWhiteArmies{
		Black: black,
		White: white,
	}, true
}

func (ship *singleShip) resetShip(shipSize, index int) {
	ship.Size = shipSize
	ship.Hits = 0
	ship.Index = index + 1
}

func (army *shipsArmy) setRandomPositions() bool {
	for tries := 0; tries < 100; tries++ {

		// reset ships & board
		shipsNumber := len(shipSizes)
		army.Ships = make([]singleShip, shipsNumber)
		army.Board = initBoard()
		for i := range army.Ships {
			army.Ships[i].resetShip(shipSizes[i], i)
		}
		success := false

		// find all ships position
		for j := range army.Ships {
			success = army.Ships[j].findPosition(army)
			if !success {
				break
			}
		}

		if success {
			return true
		}
	}

	return false
}

func (ship *singleShip) findPosition(army *shipsArmy) bool {
	found := false
	for tries := 0; tries < 100; tries++ {
		x := rand.Intn(boardSize)
		y := rand.Intn(boardSize)
		direction := rand.Intn(2)
		switch direction {
		case fromLeftToRight:
			found = ship.checkFromLeftToRight(x, y, army)
		case fromTopToBottom:
			found = ship.checkFromTopToBottom(x, y, army)
		}
		if found {
			return true
		}
	}
	return false
}

func (ship *singleShip) checkFromLeftToRight(x int, y int, army *shipsArmy) bool {
	if x+ship.Size >= boardSize {
		return false
	}
	for l := 0; l < ship.Size; l++ {
		if army.Board[x+l][y] != 0 {
			return false
		}
	}
	if x > 0 {
		army.Board[x-1][y] = freeSpaceIndex
		if y > 0 {
			army.Board[x-1][y-1] = freeSpaceIndex
		}
		if y < boardSize-1 {
			army.Board[x-1][y+1] = freeSpaceIndex
		}
	}
	if x+ship.Size < boardSize {
		army.Board[x+ship.Size][y] = freeSpaceIndex
		if y > 0 {
			army.Board[x+ship.Size][y-1] = freeSpaceIndex
		}
		if y < boardSize-1 {
			army.Board[x+ship.Size][y+1] = freeSpaceIndex
		}
	}
	for l := 0; l < ship.Size; l++ {
		if y > 0 {
			army.Board[x+l][y-1] = freeSpaceIndex
		}
		if y < boardSize-1 {
			army.Board[x+l][y+1] = freeSpaceIndex
		}
		army.Board[x+l][y] = ship.Index
	}
	return true
}

func (ship *singleShip) checkFromTopToBottom(x int, y int, army *shipsArmy) bool {
	if y+ship.Size >= boardSize {
		return false
	}
	for l := 0; l < ship.Size; l++ {
		if army.Board[x][y+l] != 0 {
			return false
		}
	}
	if y > 0 {
		army.Board[x][y-1] = freeSpaceIndex
		if x > 0 {
			army.Board[x-1][y-1] = freeSpaceIndex
		}
		if x < boardSize-1 {
			army.Board[x+1][y-1] = freeSpaceIndex
		}
	}
	if y+ship.Size < boardSize {
		army.Board[x][y+ship.Size] = freeSpaceIndex
		if x > 0 {
			army.Board[x-1][y+ship.Size] = freeSpaceIndex
		}
		if x < boardSize-1 {
			army.Board[x+1][y+ship.Size] = freeSpaceIndex
		}
	}
	for l := 0; l < ship.Size; l++ {
		if x > 0 {
			army.Board[x-1][y+l] = freeSpaceIndex
		}
		if x < boardSize-1 {
			army.Board[x+1][y+l] = freeSpaceIndex
		}
		army.Board[x][y+l] = ship.Index
	}
	return true
}
