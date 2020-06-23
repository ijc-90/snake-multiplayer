package main
import (
	constants "github.com/ijc-90/snake-multiplayer/commons"
	"fmt"
	"strings"
)

func DrawMap(aMap Map){
	topBottomLine := strings.Repeat("#", constants.Width + 2)
	fmt.Println(topBottomLine)
	for y := 0; y < constants.Height ; y ++ {
		fmt.Printf("#")
		for x := 0; x < constants.Width ; x ++ {
			if y == aMap.snakePosition.y && x == aMap.snakePosition.x{
				fmt.Printf("S")
			}else if y == aMap.fruitPosition.y && x == aMap.fruitPosition.x{
				fmt.Printf("F")
			}else{
				fmt.Printf(" ")
			}
		}

		fmt.Println("#")
	}
	fmt.Println(topBottomLine)
}