package ships

type boardArray [][]int
type shipsArmy struct {
	Ships []singleShip
	Board boardArray
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

func initBoard() boardArray {
	board := make(boardArray, boardSize)
	for i := range board {
		board[i] = make([]int, boardSize)
	}
	return board
}
