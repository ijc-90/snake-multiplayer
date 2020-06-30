package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/ijc-90/snake-multiplayer/commons"
	"image/color"
	"log"
)


var snakeOneImage *ebiten.Image
var snakeTwoImage *ebiten.Image
var fruitImage *ebiten.Image

func init() {
	var err error
	snakeOneImage, _, err = ebitenutil.NewImageFromFile("images/snake.png", ebiten.FilterDefault)
	snakeTwoImage, _, err = ebitenutil.NewImageFromFile("images/snake2.png", ebiten.FilterDefault)
	fruitImage, _, err = ebitenutil.NewImageFromFile("images/gopher2.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}
func (g *Game) Update(screen *ebiten.Image) error {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			keyPressed <- int(k)
		}
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 30*commons.Width , 30*commons.Height
}


func (g *Game) Draw(screen *ebiten.Image) {

	if !gameStarted {
		ebitenutil.DebugPrint(screen, "Loading or waiting for a match!")
	}else {
		screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

		firstPos := aMap.Snakes[0].Position.Multiply(30)
		secondPos := aMap.Snakes[1].Position.Multiply(30)
		fruitPos := aMap.FruitPosition.Multiply(30)

		var op *ebiten.DrawImageOptions
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(firstPos.X), float64(firstPos.Y))
		screen.DrawImage(snakeOneImage, op)



		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(secondPos.X), float64(secondPos.Y))
		screen.DrawImage(snakeTwoImage, op)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(fruitPos.X), float64(fruitPos.Y))
		screen.DrawImage(fruitImage, op)
	}
}