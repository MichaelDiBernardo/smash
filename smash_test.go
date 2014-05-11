package smash

import (
	"testing"
)

func TestHurt(t *testing.T) {
	sut := NewFighter(10, 0, 0, NewFixedDice([]int{}))
	sut.Hurt(5)
	if sut.HP != 5 {
		t.Errorf("Expected 5, got %d", sut.HP)
	}
}

func TestHit(t *testing.T) {
	// Attacker rolls 1, defender rolls 0.
	setD20(NewFixedDice([]int{1, 0}))

	// Attacker rolls 5 dmg.
	f1 := NewFighter(10, 4, 0, NewFixedDice([]int{5}))

	// Defender shouldn't have to roll damage.
	f2 := NewFighter(10, 0, 3, NewFixedDice([]int{}))

	f1.Attack(f2)
	if f2.HP != 5 {
		t.Errorf("Expected 5, got %d", f2.HP)
	}

	resetD20()
}

func TestMiss(t *testing.T) {
	// Attacker rolls 0, defender rolls 1.
	setD20(NewFixedDice([]int{0, 1}))

	// Attacker shouldn't have to roll damage.
	f1 := NewFighter(10, 3, 0, NewFixedDice([]int{}))

	// Defender shouldn't have to roll damage.
	f2 := NewFighter(10, 0, 3, NewFixedDice([]int{}))

	f1.Attack(f2)
	if f2.HP != 10 {
		t.Errorf("Expected 10, got %d", f2.HP)
	}

	resetD20()
}

func TestIsDead(t *testing.T) {
	sut := NewFighter(1, 0, 0, NewFixedDice([]int{}))

	if sut.Dead() {
		t.Errorf("Dead but should be alive!")
	}

	sut.Hurt(1)

	if !sut.Dead() {
		t.Errorf("Alive but should be dead!")
	}
}

func TestPickFighterFromTeam(t *testing.T) {
	team := []*Fighter{NewFighterAtRandom(), NewFighterAtRandom()}
	sut := NewTeamWithSelector(team, func(t []*Fighter) *Fighter { return t[0] })
	selected := sut.pick()

	if selected != team[0] {
		t.Errorf("Did not select dude from selector!")
	}
}

func TestTeamDead(t *testing.T) {
	team := []*Fighter{NewFighterAtRandom(), NewFighterAtRandom()}
	sut := NewTeam(team)

	if sut.Dead() {
		t.Errorf("Team isn't dead - 2 still alive!")
	}

	team[0].HP = -4

	if sut.Dead() {
		t.Errorf("Team isn't dead - 1 still alive!")
	}

	team[1].HP = -4

	if !sut.Dead() {
		t.Errorf("Team is dead - reported alive!")
	}
}

func TestFightMember(t *testing.T) {
	team := []*Fighter{
		NewFighter(5, 0, 0, NewFixedDice([]int{})),
		NewFighter(5, 0, 0, NewFixedDice([]int{})),
	}

	// No way this guy is going to miss. Will do 10 damage.
	attacker := NewFighter(10, 1000, 0, NewFixedDice([]int{10}))

	// The first guy will be selected to fight him.
	sut := NewTeamWithSelector(team, func(t []*Fighter) *Fighter { return t[0] })

	sut.DefendAgainst(attacker)

	if team[0].HP != -5 {
		t.Errorf("Expected first dude to have -5HP: had %d.", team[0].HP)
	}
}

func TestBasicBattle(t *testing.T) {
	elves := []*Fighter{
		NewFighter(5000, 5000, 5000, NewDice(4, 8)),
		NewFighter(5000, 5000, 5000, NewDice(4, 8)),
	}

	orcs := []*Fighter{
		NewFighter(1, 0, 0, NewDice(1, 1)),
		NewFighter(1, 0, 0, NewDice(1, 1)),
	}

	sut := NewBattle(NewTeam(elves), NewTeam(orcs))
	winner, _ := sut.FightItOut()

	if winner != Elves {
		t.Errorf("Elves should have won: Got %d", winner)
	}
}
