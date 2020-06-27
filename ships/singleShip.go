package ships

import "math/rand"

type singleShip struct {
	Size, Hits, Index int
}

func (ship *singleShip) resetShip(shipSize, index int) {
	ship.Size = shipSize
	ship.Hits = 0
	ship.Index = index + 1
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
