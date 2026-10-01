package main

import (
	"bufio"
	"log"
	"net"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Création, paramétrage et lancement du jeu.
func main() {

	g := game{}

	//Au lancement du jeu, le joueur se connecte au serveur.

	maConnexion, err := net.Dial("tcp", "localhost:8080")
	conn = maConnexion
	if err != nil {
		log.Println("Le serveur ne s'est pas lancé.")
		return
	}
	defer conn.Close()


	//Initialisation du Reader,du Writer et du channel.
	in = bufio.NewReader(conn)
	out = bufio.NewWriter(conn)
	myChannel = make(chan string)

	//Si le joueur == p1 (celui qui arrive en premier), g.instruction== "WAIT_OTHER_PLAYER"
	// Si le joueur ==p2 g.instruction== "BOTH_CONNECTED"

	//readFromServer lis depuis le serveur une information.
	instruction = readFromServer(conn, in)

	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}
