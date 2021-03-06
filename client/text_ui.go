package main
import (
	commons "github.com/ijc-90/snake-multiplayer/commons"
	"fmt"
	"strconv"
	"strings"
)

func DrawMap(aMap commons.Map, playerId int){
	if aMap.GameOver{
		if aMap.Snakes[playerId-1].Won{
			println("###### Congratz! you won! ######")
		}else if aMap.Snakes[playerId-1].Lost{
			println("###### Shame on you, noob! you lost! ######")
		}else{
			println("###### Its a tie! ######")
		}
		println("Game is over, press enter key to exit")
		return
	}
	topBottomLine := strings.Repeat("#", (aMap.Width + 2)*3)
	fmt.Println(topBottomLine)

	for y := 0; y < aMap.Height ; y ++ {

		fmt.Printf("###")
		for x := 0; x < aMap.Width ; x ++ {
			if y == aMap.Snakes[0].Position.Y && x == aMap.Snakes[0].Position.X{
				snakeId := strconv.Itoa(aMap.Snakes[0].Id)
				fmt.Printf("S%sS", snakeId)
			}else if y == aMap.Snakes[1].Position.Y && x == aMap.Snakes[1].Position.X{
				snakeId := strconv.Itoa(aMap.Snakes[1].Id)
				fmt.Printf("S%sS", snakeId)
			}else if y == aMap.FruitPosition.Y && x == aMap.FruitPosition.X{
				fmt.Printf("FFF")
			}else{
				fmt.Printf("   ")
			}
		}

		fmt.Println("###")
	}
	fmt.Println(topBottomLine)
}