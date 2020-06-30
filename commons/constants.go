package commons

const Width = 20
const Height = 15

var Directions = map[rune]int{
	'w' : 1, // UP
	32 : 1, // UP
	'a' : 2, // LEFT
	10 : 2, // LEFT
	's' : 3, // DOWN
	28 : 3, // DOWN
	'd' : 4, // RIGHT
	13 : 4, // RIGHT

}

const TickInterval = 500
