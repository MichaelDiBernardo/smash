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

    // Defender shouldn't have to roll anything.
    f2 := NewFighter(10, 0, 3, NewFixedDice([]int{}))

    f1.Attack(f2)
    if f2.HP != 5 {
		t.Errorf("Expected 5, got %d", f2.HP)
    }
}
