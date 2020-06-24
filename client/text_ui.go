package main
import (
	commons "github.com/ijc-90/snake-multiplayer/commons"
	"fmt"
	"strings"
)

func DrawMap(aMap commons.Map){
	topBottomLine := strings.Repeat("#", (aMap.Width + 2)*3)
	fmt.Printf(topBottomLine)
	fmt.Printf("\n")
	for y := 0; y < aMap.Height ; y ++ {
		fmt.Printf("###")
		for x := 0; x < aMap.Width ; x ++ {
			if y == aMap.SnakePosition.Y && x == aMap.SnakePosition.X{
				fmt.Printf("SSS")
			}else if y == aMap.FruitPosition.Y && x == aMap.FruitPosition.X{
				fmt.Printf("FFF")
			}else{
				fmt.Printf("   ")
			}
		}

		fmt.Println("###")
		fmt.Printf("\n")
	}
	fmt.Println(topBottomLine)
	fmt.Printf("\n")
}