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
	D20 = NewFixedDice([]int{1, 0})

	// Attacker rolls 5 dmg.
	f1 := NewFighter(10, 4, 0, NewFixedDice([]int{5}))

	// Defender shouldn't have to roll damage.
	f2 := NewFighter(10, 0, 3, NewFixedDice([]int{}))

	f1.Attack(f2)
	if f2.HP != 5 {
		t.Errorf("Expected 5, got %d", f2.HP)
	}
}

func TestMiss(t *testing.T) {
	// Attacker rolls 0, defender rolls 1.
	D20 = NewFixedDice([]int{0, 1})

	// Attacker shouldn't have to roll damage.
	f1 := NewFighter(10, 3, 0, NewFixedDice([]int{}))

	// Defender shouldn't have to roll damage.
	f2 := NewFighter(10, 0, 3, NewFixedDice([]int{}))

	f1.Attack(f2)
	if f2.HP != 10 {
		t.Errorf("Expected 10, got %d", f2.HP)
	}
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

func TestSelectFighterFromTeam(t *testing.T) {
	team := []*Fighter{NewFighterAtRandom(), NewFighterAtRandom()}
	sut := NewTeamWithSelector(team, func(t []*Fighter) *Fighter { return t[0] })
    selected := sut.Select()

	if selected != team[0] {
		t.Errorf("Did not select dude from selector!")
	}
}
