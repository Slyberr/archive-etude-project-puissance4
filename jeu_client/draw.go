package main

import (
	"image/color"
	"strconv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g *game) Draw(screen *ebiten.Image) {

	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case titleState:
		g.titleDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}

}

// Affichage des graphismes de l'écran titre.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)

	if g.stateFrame >= globalBlinkDuration/3 {

		//Si le joueur appuie en premier, il recevera sur son écran une instruction disant d'attendre l'autre joueur.
		if bePress >= 1 {
			text.Draw(screen, "En attente de votre Adversaire...", smallFont, 110, 500, globalTextColor)

		} else {
			text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
		}
	}

}

// Affichage des graphismes de l'écran de sélection des couleurs des joueurs.
func (g game) colorSelectDraw(screen *ebiten.Image) {
	
	//Trois conditions: 
	
	//Soit bePress vaut au moins 1, et dans cela signifie que la couleur à été selectionnée DÉFINITEVEMENT 
	if  bePress >= 1{
		text.Draw(screen, "En attente de votre Adversaire...\n n°Couleur choisie:"+strconv.Itoa(g.p1Color), smallFont, 60, 100, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	
	//Soit les deux clients sont sur la même couleur.
	} else if warningSameColor{
		text.Draw(screen, "Attention, Même couleur que l'adversaire!\n(pas selectionnable = déjà prise...)", smallFont, 60, 100,globalTextColor)
	//Cas de base.
	} else {
		text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 100, globalTextColor)
	}

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line

		//Si jamais p1Color=p2Color, alors p1 aura un cercle avec un diamètre plus grand pour bien voir la différence entre les deux cercles.
		if numColor == g.p1Color {
			if g.p1Color == g.p2Color {
				vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2+5, globalSelectColor, true)

			} else {
				vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColor, true)

			}
		}
		//Pour mettre à bien l'extension du curseur partagé, il faut pouvoir afficher le curseur de l'adversaire.
		if numColor == g.p2Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, true)

		}
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)
	vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)

	//Affiche si le joueur peut jouer ou non.
	if g.turn == p1Turn {
		text.Draw(screen, "A Votre tour !", smallFont, 10, 20, globalTextColor)

	} else {
		text.Draw(screen, "Au tour de votre adversaire...", smallFont, 20, 40, globalTextColor)
	}
 
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}

	//Dès que le joueur veut rejouer (bePress >=1), alors il y a un message en attendant que l'adversaire accepte également de jouer.
	if bePress >= 1 {
		text.Draw(screen, "Vous voulez rejouer ? \nVeuillez attendre la réponse de votre adversaire.", smallFont, 0, 400, globalTextColor)
	} else {
		text.Draw(screen, message, smallFont, 300, 400, globalTextColor)
	}
}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
