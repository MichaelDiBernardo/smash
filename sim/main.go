package main

import (
	"fmt"
	"math/rand"
	"time"

	"smash"
)

// Create a random team of dudes. Each dude created will send 1 to the given
// channel updon dying.
func randomTeam(deathQueue chan int) *smash.Team {
	dudes := []*smash.Fighter{
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
		smash.NewFighterAtRandom(),
	}

	onDeath := func() {
		deathQueue <- 1
	}

	for _, dude := range dudes {
		dude.OnDeath = onDeath
	}

	return smash.NewTeam(dudes)
}

// Attempt to pull numDeaths values off of the given deathQueue channel. Will
// send the given allegiance to the doneSignal channel when this happens,
// signalling that this allegiance won.
func countDeaths(deathQueue chan int, doneSignal chan int, numDeaths int, allegiance int) {
	for i := 0; i < numDeaths; i++ {
		<-deathQueue
	}
	doneSignal <- allegiance
}

// Fight it out. Get the winners and throw them on the winner queue for their
// allegiance, so that they can be assigned to another battle.
func doBattle(battle *smash.Battle, winnerQueues []chan *smash.Team) {
	allegiance, winners := battle.FightItOut()
	winnerQueues[allegiance] <- winners
}

// Watch the winner queues for both allegiances and when you have a pair of
// winning teams, throw them at one another.
func matchMaker(winnerQueues []chan *smash.Team) {
	for {
		winningElves := <-winnerQueues[smash.Elves]
		winningOrcs := <-winnerQueues[smash.Orcs]

		battle := smash.NewBattle(winningElves, winningOrcs)
		go doBattle(battle, winnerQueues)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// How many dudes do we have participating?
	nTeams := 10
	nFightersPerSide := nTeams * 3

	// Matchmaking setup.
	winnerQueues := []chan *smash.Team{
		make(chan *smash.Team),
		make(chan *smash.Team),
	}
	go matchMaker(winnerQueues)

	// Win condition setup. (Once everyone of an allegiance dies, we know the
	// other allegiance won.
	elfDeaths := make(chan int)
	orcDeaths := make(chan int)
	doneSignal := make(chan int)

	// Go and count elf deaths; if they all die, orcs won.
	go countDeaths(elfDeaths, doneSignal, nFightersPerSide, smash.Orcs)
	// And vice-versa.
	go countDeaths(orcDeaths, doneSignal, nFightersPerSide, smash.Elves)

	// Start the fights.
	for i := 0; i < 10; i++ {
		battle := smash.NewBattle(randomTeam(elfDeaths), randomTeam(orcDeaths))
		// Best possible use of the 'go' keyword.
		go doBattle(battle, winnerQueues)
	}

	// Wait one of the countDeaths goroutines to tell us that there was a winner.
	winner := <-doneSignal
	fmt.Printf("The winner is %d!\n", winner)
}
