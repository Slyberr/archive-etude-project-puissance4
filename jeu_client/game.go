package main

import (
	"bufio"
	"net"
	"sync"
)

// Structure de données pour représenter l'état courant du jeu.
type game struct {
	gameState     int
	stateFrame    int
	grid          [globalNumTilesX][globalNumTilesY]int
	p1Color       int
	p2Color       int
	turn          int
	tokenPosition int
	result        int
}

// Constantes pour représenter la séquence de jeu actuelle (écran titre,
// écran de sélection des couleurs, jeu, écran de résultats).
const (
	titleState int = iota
	colorSelectState
	playState
	resultState
)

// Constantes pour représenter les pions dans la grille de puissance 4
// (absence de pion, pion du joueur 1, pion du joueur 2).
const (
	noToken int = iota
	p1Token
	p2Token
)

// Constantes pour représenter le tour de jeu (joueur 1 ou joueur 2).
const (
	p1Turn int = iota
	p2Turn
)

// Constantes pour représenter le résultat d'une partie (égalité si
// la grille a été remplie sans qu'un joueur n'ait gagné, joueur 1
// gagnant ou joueur 2 gagnant).
const (
	equality int = iota
	p1wins
	p2wins
)

// Une connection avec le seveur
var conn net.Conn

// Valeur en chaîne de caractère de l'information envoyée par un client, passe par le serveur qui la renvoie ensuite  à l'autre client.
var instruction string

// récupère des informations, toujours sous forme de string.
var in *bufio.Reader

// envoie des informations
var out *bufio.Writer

// Un channel pour pouvoir utiliser nos go routines correctement.
var myChannel chan string

// "bePress"== a pressé la touche qui permet de passer à la vue suivante.
// Si le joueur n'a pas pressé la dite touche, bePress=0. S'il la pressée, bePress=1
// Une valeur numérique permet, à la place d'un booléen, de contrôler si le joueur appuye plusieurs fois sur la touche entrée pour une même vue.
// Le programme a été conçu pour que l'information ne soit envoyée qu'une seule fois (et donc éviter de surcharger le serveur). 
// Cela permettera donc de ne pas réenvoyer la même information au serveur.
var bePress int

// Verifie que la go routine ne soit pas déjà appellé pour éviter de surcharger de go routine pour rien (60 fois par seconde dans Update()).
var goRoutineWasCall bool

// = true si la couleur adverse a été selectionnée.
var adverseColorSelected bool

//Vaut vrai lorsque le curseur du joueur adverse se situe au même endroit que le joueur courant.
var warningSameColor bool


//Permet d'assurer l'intégrité des fonctions utilisées, ici elle est uniquement utilisée pour readFromServerForRoutine().
var monVerrou sync.Mutex


// Remise à 0 du jeu pour recommencer une partie. Le joueur qui a
// perdu la dernière partie commence.
func (g *game) reset() {
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}
}
