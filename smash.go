package smash

import (
	"math/rand"
)

//------------------------------------------------------------------------------
// Random numbers.
//------------------------------------------------------------------------------

// Anything that can produce a random value is a roller. Bad name, I know.
type Roller interface {
	Roll() int
}

// Concrete random dice.
type Dice struct {
	die   int
	sides int
}

func NewDice(die int, sides int) *Dice {
	return &Dice{die: die, sides: sides}
}

// Assemble dice with 1-3 die and 4-7 sides.
func NewDiceAtRandom() *Dice {
	die := rand.Intn(3) + 1   // 1 to 3 die
	sides := rand.Intn(4) + 4 // 4 to 7 sides.
	return NewDice(die, sides)
}

func (self *Dice) Roll() int {
	total := 0
	for i := 1; i <= self.die; i++ {
		total += i * (rand.Intn(self.sides) + 1)
	}
	return total
}

// Fixed dice that get values from a list.
type FixedDice struct {
	vals []int
	cur  int
}

func NewFixedDice(vals []int) *FixedDice {
	return &FixedDice{vals: vals, cur: 0}
}

func (self *FixedDice) Roll() int {
	val := self.vals[self.cur]
	self.cur = (self.cur + 1) % len(self.vals)
	return val
}

//------------------------------------------------------------------------------
// Fighters and their stuff.
//------------------------------------------------------------------------------
var D20 Roller = NewDice(1, 20)

// A dude who fights.
type Fighter struct {
	HP      int
	melee   int
	evasion int
	dice    Roller
}

func NewFighter(hp int, melee int, evasion int, dice Roller) *Fighter {
	return &Fighter{HP: hp, melee: melee, evasion: evasion, dice: dice}
}

func NewFighterAtRandom() *Fighter {
	hitDice := NewDice(6, 4)
	skillDice := NewDice(10, 2)
	return NewFighter(hitDice.Roll(), skillDice.Roll(), skillDice.Roll(), NewDiceAtRandom())
}

func (self *Fighter) Hurt(dmg int) {
	self.HP -= dmg
}

func (self *Fighter) Attack(other *Fighter) {
	atk := D20.Roll() + self.melee
	ev := D20.Roll() + other.evasion

	if ev > atk {
		return
	}

	other.Hurt(self.dice.Roll())
}

func (self *Fighter) Dead() bool {
    return self.HP <= 0
}

//------------------------------------------------------------------------------
// Battles.
//------------------------------------------------------------------------------
