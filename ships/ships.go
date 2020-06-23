package ships

import (
	"encoding/gob"
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
type hitReport struct {
	Hit, Sank bool
	X, Y      int
}

func Init(writer http.ResponseWriter, request *http.Request) bool {
	gob.Register(BlackAndWhiteArmies{})

	session := sessions.CheckSession(writer, request)
	ships, success := initShips()
	if !success {
		return false
	}

	session.Values["ships"] = ships
	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
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

// This is brute-force algorithm. For production I would use something more clever.
// There are 100 tries to create a board.
// Each ship's position is resolved in less than 100 tries.
// Ships positions are saved in shipsArmy.Board two dimensional slice.
// First dimension is horizontal, second dimension is vertical.
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

		// find all ships' position
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

// try up to 100 times until you find ship's position
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

// There is only checking from left to right, not right to left, as the result is the same.
// Please notice how freeSpaceIndex is assigned to space around the ship.
// This makes the space occupied, so no other ship can take it.
func (ship *singleShip) checkFromLeftToRight(x int, y int, army *shipsArmy) bool {
	// if ship cannot fit inside the board at this position
	if x+ship.Size >= boardSize {
		return false
	}

	// find if space for ship is not occupied already
	for l := 0; l < ship.Size; l++ {
		if army.Board[x+l][y] != 0 {
			return false
		}
	}
	// space for the ship is found at this point

	// reserve space at the left side of the ship
	if x > 0 {
		army.Board[x-1][y] = freeSpaceIndex
		if y > 0 {
			army.Board[x-1][y-1] = freeSpaceIndex
		}
		if y < boardSize-1 {
			army.Board[x-1][y+1] = freeSpaceIndex
		}
	}
	// reserve space at the right side of the ship
	if x+ship.Size < boardSize {
		army.Board[x+ship.Size][y] = freeSpaceIndex
		if y > 0 {
			army.Board[x+ship.Size][y-1] = freeSpaceIndex
		}
		if y < boardSize-1 {
			army.Board[x+ship.Size][y+1] = freeSpaceIndex
		}
	}
	// reserve space at the top and bottom of the ship and the ship itself
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

// There is only checking from top to bottom, not bottom to top, as the result is the same.
// Please notice how freeSpaceIndex is assigned to space around the ship.
// This makes the space occupied, so no other ship can take it.
func (ship *singleShip) checkFromTopToBottom(x int, y int, army *shipsArmy) bool {
	// if ship cannot fit inside the board at this position
	if y+ship.Size >= boardSize {
		return false
	}

	// find if space for ship is not occupied already
	for l := 0; l < ship.Size; l++ {
		if army.Board[x][y+l] != 0 {
			return false
		}
	}
	// space for the ship is found at this point

	// reserve space at the top of the ship
	if y > 0 {
		army.Board[x][y-1] = freeSpaceIndex
		if x > 0 {
			army.Board[x-1][y-1] = freeSpaceIndex
		}
		if x < boardSize-1 {
			army.Board[x+1][y-1] = freeSpaceIndex
		}
	}
	// reserve space at the bottom of the ship
	if y+ship.Size < boardSize {
		army.Board[x][y+ship.Size] = freeSpaceIndex
		if x > 0 {
			army.Board[x-1][y+ship.Size] = freeSpaceIndex
		}
		if x < boardSize-1 {
			army.Board[x+1][y+ship.Size] = freeSpaceIndex
		}
	}
	// reserve space at the left and right side of the ship and the ship itself
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

func (armies *BlackAndWhiteArmies) Fire(whoFires string, x, y int) hitReport {
	var armyDefending shipsArmy

	switch whoFires {
	case "white":
		armyDefending = armies.Black
	case "black":
		armyDefending = armies.White
	}

	// and shoot
	hit := false
	sank := false
	index := armyDefending.Board[x][y]
	if index > 0 && index < freeSpaceIndex {
		hit = true
		shipIndex := index - 1
		armyDefending.Ships[shipIndex].Hits++
		if armyDefending.Ships[shipIndex].Hits >= armyDefending.Ships[shipIndex].Size {
			sank = true
		}
	}

	return hitReport{
		Hit:  hit,
		Sank: sank,
		X:    x,
		Y:    y,
	}
}

func GetTarget() (x, y int) {
	return rand.Intn(boardSize), rand.Intn(boardSize)
}
