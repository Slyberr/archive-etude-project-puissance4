package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Constantes définissant les paramètres généraux du programme.
const (
	globalWidth         = globalNumTilesX * globalTileSize
	globalHeight        = (globalNumTilesY + 1) * globalTileSize
	globalTileSize      = 100
	globalNumTilesX     = 7
	globalNumTilesY     = 6
	globalCircleMargin  = 5
	globalBlinkDuration = 60
	globalNumColorLine  = 3
	globalNumColorCol   = 3
	globalNumColor      = globalNumColorLine * globalNumColorCol
)

// Variables définissant les paramètres généraux du programme.
var (
	//nous avons changé quelque scouleurs qui nous conviennent +
	globalBackgroundColor color.Color = color.NRGBA{R: 176, G: 196, B: 222, A: 255}
	globalGridColor       color.Color = color.NRGBA{R: 150, G: 150, B: 170, A: 255}
	globalTextColor       color.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	globalSelectColor     color.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	smallFont             font.Face
	largeFont             font.Face
	globalTokenColors     [globalNumColor]color.Color = [globalNumColor]color.Color{
		color.NRGBA{R: 255, G: 0, B: 255, A: 255},
		color.NRGBA{R: 0, G: 255, B: 255, A: 255},
		color.NRGBA{R: 255, G: 255, B: 0, A: 255},
		color.NRGBA{R: 0, G: 255, B: 0, A: 255},
		color.NRGBA{R: 255, G: 128, B: 0, A: 255},
		color.NRGBA{R: 0, G: 128, B: 255, A: 255},
		color.NRGBA{R: 255, G: 0, B: 0, A: 255},
		color.NRGBA{R: 221, G: 160, B: 221, A: 255},
		color.NRGBA{R: 88, G: 44, B: 0, A: 255},
	}
	offScreenImage *ebiten.Image
)
