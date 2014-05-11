package smash

import (
	"fmt"
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
		total += rand.Intn(self.sides) + 1
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
var defaultD20 Roller = NewDice(1, 20)
var d20 Roller = defaultD20

func setD20(r Roller) {
	d20 = r
}

func resetD20() {
	d20 = defaultD20
}

// A dude who fights.
type Fighter struct {
	HP      int
	melee   int
	evasion int
	dice    Roller
	OnDeath func()
}

func NewFighter(hp int, melee int, evasion int, dice Roller) *Fighter {
	return &Fighter{
		HP:      hp,
		melee:   melee,
		evasion: evasion,
		dice:    dice,
		OnDeath: func() {},
	}
}

func NewFighterAtRandom() *Fighter {
	hitDice := NewDice(6, 4)
	skillDice := NewDice(10, 2)
	return NewFighter(
		hitDice.Roll(),
		skillDice.Roll(),
		skillDice.Roll(),
		NewDiceAtRandom(),
	)
}

func (self *Fighter) Hurt(dmg int) {
	if self.Dead() {
		panic("Hitting a dead man!")
	}

	self.HP -= dmg
	fmt.Printf("%d dmg. CurHP: %d\n", dmg, self.HP)

	if self.Dead() {
		self.OnDeath()
	}
}

func (self *Fighter) Attack(other *Fighter) {
	atk := d20.Roll() + self.melee
	ev := d20.Roll() + other.evasion

	fmt.Printf("Rolled %d against %d.\n", atk, ev)

	if ev > atk {
		return
	}

	other.Hurt(self.dice.Roll())
}

func (self *Fighter) Dead() bool {
	return self.HP <= 0
}

//------------------------------------------------------------------------------
// Teams.
//------------------------------------------------------------------------------
type Team struct {
	roster   []*Fighter
	selector func([]*Fighter) *Fighter
}

func NewTeamWithSelector(roster []*Fighter, selector func([]*Fighter) *Fighter) *Team {
	return &Team{roster: roster, selector: selector}
}

func defaultSelector(roster []*Fighter) *Fighter {
	n := rand.Intn(len(roster))
	f := roster[n]
	if f.Dead() {
		return defaultSelector(roster)
	}
	return f
}

func NewTeam(roster []*Fighter) *Team {
	return NewTeamWithSelector(roster, defaultSelector)
}

func (self *Team) Dead() bool {
	for _, f := range self.roster {
		if !f.Dead() {
			return false
		}
	}
	return true
}

func (self *Team) DefendAgainst(other *Fighter) {
	defender := self.pick()
	other.Attack(defender)
}

func (self *Team) pick() *Fighter {
	return self.selector(self.roster)
}

//------------------------------------------------------------------------------
// Battles.
//------------------------------------------------------------------------------
const (
	Elves = iota
	Orcs
)

type Battle struct {
	teams []*Team
}

func NewBattle(elves *Team, orcs *Team) *Battle {
	teams := []*Team{elves, orcs}
	return &Battle{teams: teams}
}

func (self *Battle) FightItOut() (int, *Team) {
	// Select random allegiance to start.
	atkInd := rand.Intn(2)

	for {
		atkInd = (atkInd + 1) % 2
		defInd := (atkInd + 1) % 2

		fmt.Printf("Atk: %d Def: %d\n", atkInd, defInd)

		attacker := self.teams[atkInd]
		defender := self.teams[defInd]

		champion := attacker.pick()
		defender.DefendAgainst(champion)

		if defender.Dead() {
			return atkInd, self.teams[atkInd]
		}
	}
}
