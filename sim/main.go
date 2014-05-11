package main

import (
    "fmt"
    "math/rand"
    "time"

	"smash"
)

func main() {
    rand.Seed(time.Now().UnixNano())

	elves := smash.NewTeam([]*smash.Fighter{
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
	})

	orcs := smash.NewTeam([]*smash.Fighter{
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
	})

    battle := smash.NewBattle(elves, orcs)
    winner := battle.FightItOut()
    fmt.Printf("Winner is: %d\n", winner)
}
