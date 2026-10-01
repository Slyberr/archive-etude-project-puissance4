package main

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	//Cas vue écran Titre.
	case titleState:

		//bePress permet de vérifier si la touche entrée a été pressée précédemment (dans une même vue).
		//Si oui, alors il ne rentrera plus dans cette condition. (Détails pour cette variable dans le fichier game.go)
		if g.titleUpdate() && bePress == 0 {

			bePress++
			//Si le client appuie sur entrée, alors il devra attendre l'autre client pour que la vue suivante s'affiche.
			writeForServer(conn, out, "OK_COLOR")
			go readFromServerForRoutine(conn, in, myChannel)
		}

		//Le select permettant de lire l'information de la goroutine est en dehors de g.titleUpdate(), car on passe dans cette fonction uniquement si la touche entrée a bien été pressée.
		//Nous voulons que le passage à la vue suivante soit gérée indépendemment de cela.
		select {
		//Les deux joueurs devront passer par là, si l'instruction lue depuis le serveur contient bien "OK_COLOR", alors ils peuvent passer à la vue suivante.
		//(Les deux s'envoient et doivent recevoir "OK_COLOR")
		case instruction = <-myChannel:

			if instruction == "OK_COLOR" {
				//On remet le bePress à 0 pour pouvoir le réutiliser pour la suite.
				bePress = 0
				g.gameState++
			}
		default:
			//Si rien n'est lu, alors Update() Continue de s'exécuter normalement.
		}
	//Cas Vue Selection de couleur
	case colorSelectState:

		//Dans tous les cas , on lis l'information venant du server (car même si notre couleur à été selectionnée,
		//on voudrait voir les déplacements du curseur du joueur adverse !)
		go readFromServerForRoutine(conn, in, myChannel)

		//Nous souhaitons pouvoir faire une distinction entre mise à jour du curseur adverse et selection d'une couleur:
		if bePress == 0 {
			//Si jamais la touche entrée n'a jamais été pressée:

			//Si colorSelectUpdate==true, alors entrée a été pressé.
			if g.colorSelectUpdate() {

				//Lorsque la couleur est sélectionnée, bePress++.
				//Donc ne passera plus là <=> LE curseur reste bloqué et ne peut plus être mis à jour.
				bePress++

				//On informe le joueur adverse avec "SELECTED" à la fin de l'envoi de notre couleur.
				writeForServer(conn, out, strconv.Itoa(g.p1Color)+"SELECTED")

			} else {

				//Si jamais la couleur reçue de l'adversaire est la même que celle du client,
				//Alors on warningSameColor=true, et cela nous sert de condition pour mettre à jour le texte
				//de la vue.
				if string(instruction[:1]) == strconv.Itoa(g.p1Color) {
					warningSameColor = true
				} else {
					warningSameColor = false
				}

				//Si jamais ça n'a pas été pressé, étant donné qu'on désire voir les déplacements de curseurs de l'adversaire,
				//On doit tout de même envoyer notre position actuelle.
				writeForServer(conn, out, strconv.Itoa(g.p1Color))
			}
		}

		select {
		case instruction = <-myChannel:

			//Si l'info obtenue est un simple caratère, alors c'est juste une mise à jour de curseur.
			if len(instruction) == 1 {
				g.p2Color, _ = strconv.Atoi(instruction)
				//La couleur de l'adversaire n'est donc pas selectionnée.
				adverseColorSelected = false

				//COMMENTAIRE A REFORMULER.
				//On va vouloir lire "UR_PLAYER1" ou "UR_PLAYER2" (Info envoyée par le serveur).
				//On sera capable de le lire à la ligne 49.
				//(qui nous sert donc à la fois de lecture de curseur de couleur/Selection de couleur et de lecture attribution du joueur).
				//On verifira ensuite la lecture par la condition à la ligne 95.
			} else if (instruction == "UR_PLAYER1") || (instruction == "UR_PLAYER2") {
				bePress = 0
				if instruction == "UR_PLAYER2" {
					g.turn = p2Turn
				}
				g.gameState++

				//Si l'info obtenue n'est ni un simple caractère, ni une info (envoyée par le serveur) exprimant qu'on est joueur 1 ou 2,
				//Alors il s'agit de l'info où le joueur adverse a bien selectionné sa couleur.
				//(Dans le serveur cela est symbolisée par une condition len(instruction >1) )
			} else {
				g.p2Color, _ = strconv.Atoi(instruction[:1])
				//L'adversaire a bien selectionné définitivement la couleur, donc adversaireColorSelected = true.
				adverseColorSelected = true
			}

		default:

		}
	//Vue du jeu.
	case playState:
		//Si une touche droite ou gauche est pressé, ça change la position du jeton.
		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int

		if g.turn == p1Turn {

			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
			if lastXPositionPlayed != -1 && lastYPositionPlayed != -1 {
				//Si le jeton est joué, alors on va envoyer les coordonnées du jeton joué à l'adversaire pour que sa vue soit mise à jour.
				bePress++
				//On envoie la coordonnée X et Y du pion joué.
				chaineAEnvoyer := strconv.Itoa(lastXPositionPlayed) + strconv.Itoa(lastYPositionPlayed)
				writeForServer(conn, out, chaineAEnvoyer)
			}

		} else {

			//Les coordonnées du jeton de l'adversaire sont prises en compte, et sont mise à jour dans la grid qui elle même mettera à jour la vue.

			//Permet d'appeller une seule fois la goRoutine.
			if !goRoutineWasCall {
				go readFromServerForRoutine(conn, in, myChannel)
				goRoutineWasCall = true
			}
		}
		select {
		case instruction = <-myChannel:
			//Lorsque l'instruction est reçue, monX est égal au premier caractère de la chaîne envoyé par le client.
			monX, _ := strconv.Atoi(string(instruction[0]))
			bePress++
			//p2Update est maintenant une fonction qui prend en argument la coordonnée X étant donné qu'on a l'information venant de l'extérieur
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update(monX)
			//On met goRoutinewasCall à false car dès que l'information est reçue, on sera par la suite
			//Re apte à recevoir à nouveau l'information de la position d'une autre jeton de l'adversaire.
			goRoutineWasCall = false
		default:

		}

		//Le bePress ici nous assure que le jeton à été joué pour le client.
		//Dans le cas où c'est l'adversaire qui joue, on regardera si le jeu se termine uniquement si l'autre joueur à bien joué (bePress++ à la ligne 126).
		if lastXPositionPlayed >= 0 && bePress >= 1 {

			bePress = 0
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			//Dès qu'on check si le jeu se termine, on remet bien bePress à 0 pour éviter de passer par là alors que rien ne se passe dans le jeu.
			if finished {
				g.result = result
				//Les informations envoyé sur le résultat seront interprété par le serveur.
				//(car si le joueur perd, il commence en premier.)
				if g.result == p1wins {
					writeForServer(conn, out, "END_GAME_I_WIN")
				} else if g.result ==p2wins{
					writeForServer(conn, out, "END_GAME_I_LOSE")
				}else{
					writeForServer(conn,out,"END_GAME_EQUALITY")
				}
				g.gameState++
			}

		}

	//Vue où on voit le résultat.
	case resultState:

		//On dit à l'adversaire qu'on est prêt pour en refaire une, on attend sa confirmation.
		//bePress a encore le même principe pour éviter d'envoyer la même information.

		//En réalité, on suite exactement la même méthode que pour "OK_COLOR" dans la vue titleState.

		if g.resultUpdate() && bePress == 0 {
			bePress++
			writeForServer(conn, out, "ANOTHER_GAME_PLS")
			go readFromServerForRoutine(conn, in, myChannel)

		}

		select {
		case instruction = <-myChannel:
			//On oublie pas de reset le jeu.
			if instruction == "ANOTHER_GAME_PLS" {
				bePress = 0
				g.reset()
				g.gameState = playState
			}
		default:
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		col = (col + 1) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		line = (line + 1) % globalNumColorLine
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
	}

	g.p1Color = line*globalNumColorLine + col

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {

		//On souhaite renvoyer faux uniquement si l'adversaire à déjà selectionné sa couleur. Pourquoi ?
		//Car on veut pouvoir faire en sorte que si les deux joueurs ont le curseur sur la même couleur mais que aucun n'a encore selectionné,
		//Il est possible de selectionner la couleur pour le premier.
		if g.p2Color == g.p1Color && adverseColorSelected {
			return false
		}

		return true
	}

	return false
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// La fonction est grandement changée, elle prend en argument un seul int (le X) et n'a pas besoin de vérifier si la position pose problème (ligne remplie).
// En effet, elle a été validée par le joueur adverse, qui lui a testé avant si cela était envisageable.
// De plus elle return les mêmes arguments que p1, pour pouvoir vérifier si le jeu se termine.
func (g *game) p2Update(monX int) (int, int) {

	_, yPos := g.updateGrid(p2Token, monX)
	lastXPositionPlayed := monX
	lastYPositionPlayed := yPos

	g.turn = p1Turn
	return lastXPositionPlayed, lastYPositionPlayed
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g game) resultUpdate() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {

		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {

		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {

		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
