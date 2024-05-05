package main

import (
	"math/rand"
)

type Dice uint8

// Random return a number between 1 and 6 included
func random6() uint8 {
	return uint8(rand.Intn(6) + 1)
}

// rollDices takes the number i of dices to roll and returns a slice of i pseudo random numbers between 1 and 6 included.
func RollDices(i int) []Dice { 
 	res := make([]Dice, i)
	for c, _ := range res {
		res[c] = Dice(random6())
	}
	return res
}