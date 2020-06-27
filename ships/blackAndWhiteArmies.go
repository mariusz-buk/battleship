package ships

type BlackAndWhiteArmies struct {
	White, Black shipsArmy
}

type hitReport struct {
	Hit, Sank bool
	X, Y      int
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
