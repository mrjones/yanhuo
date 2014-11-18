package main

import (
	"fmt"
)

type color int8
type value int8

const (
	WHITE color = iota
	RED
	BLUE
	YELLOW
	GREEN

)

var ALL_COLORS = []color{WHITE, RED, BLUE, YELLOW, GREEN}

type colorInfo struct {
	fullName string
	shortName string
}

var COLOR_INFOS = map[color]colorInfo {
	WHITE: colorInfo{fullName: "WHITE", shortName: "W"},
	RED: colorInfo{fullName: "RED", shortName: "R"},
	BLUE: colorInfo{fullName: "BLUE", shortName: "B"},
	YELLOW: colorInfo{fullName: "YELLOW", shortName: "Y"},
	GREEN: colorInfo{fullName: "GREEN", shortName: "G"},
}


var VALUE_COUNTS = map[value]int{
	1: 3,
	2: 2,
	3: 2,
	4: 2,
	5: 1,
}

func main() {
	fmt.Println("Hello, world!")
	
	_ = createDeck()
}



type card struct {
	value value
	color color
}

func createDeck() []card {
	cards := []card{}
	for color, _ := range(COLOR_INFOS) {
		for value, count := range(VALUE_COUNTS) {
			for i := 0; i < count; i++ {
				cards = append(cards, card{value: value, color: color})
				fmt.Printf("color: %s, value: %d\n", COLOR_INFOS[color].fullName, value)
			}
		}
	}
	
	return cards
}
