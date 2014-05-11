package main

import (
    "math/rand"
    "time"

	"smash"
)

func randomTeam() *smash.Team {
	return smash.NewTeam([]*smash.Fighter{
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
	})
}

func doBattle(battle *smash.Battle, winnerQueues []chan *smash.Team) {
    allegiance, winners := battle.FightItOut()
    winnerQueues[allegiance] <- winners
}

func main() {
    rand.Seed(time.Now().UnixNano())

    winnerQueues := []chan *smash.Team{
        make(chan *smash.Team),
        make(chan *smash.Team),
    }

    for i := 0; i < 10; i++ {
        battle := smash.NewBattle(randomTeam(), randomTeam())
        doBattle(battle, winnerQueues)
    }
}
