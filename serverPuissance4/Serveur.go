package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {

	//Capable d'écouter
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	//Première connexion + instanciation Reader et Writer.
	//Le premier client se verra attribué le statut "WAIT_OTHER_PLAYER"
	conn, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	defer conn.Close()
	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)

	writeForClient(conn, out, "WAIT_OTHER_PLAYER")


	//Acceptation de la seconde connexion et IDEM.
	//Le second client se verra attribué le statut "BOTH_CONNECTED"

	conn2, err2 := listener.Accept()
	if err2 != nil {
		log.Println("accept error:", err2)
		return
	}

	in2 := bufio.NewReader(conn2)
	out2 := bufio.NewWriter(conn2)

	writeForClient(conn2, out2, "BOTH_CONNECTED")
	defer conn2.Close()


	//Les joueurs vont attendre de l'autre qu'il reçoit "OK_COLOR".
	//Le premier qui recevera cela vera "Attente du second joueur", tandis que le second qui appuie déclenchera la vue de l'autre ainsi que la sienne.
	okVueColor := readFromClient(conn2, in2)
	okVueColor2 := readFromClient(conn, in)

	writeForClient(conn2, out2, okVueColor2)
	writeForClient(conn, out, okVueColor)

	//Protocole dans lequel Le client 1 le curseur de sa couleur au client 2 et inversement.
	selectCircle1 := ""
	selectCircle2 := ""

	compteurBreakLoop := 0

	//A chaque tour de boucle, on veut pouvoir recevoir et envoyer l'info du curseur du client.
	//Lorsque un des selectcircle a une longueur >1, cela signifie que la string "SELECTED" se trouve derrière.
	//Dans ce cas là, On ne rentrera plus jamais dans la condition (De plus le curseur est bloquée du côté client, donc plus besoin de renvoyer l'information).
	//Si les deux conditions ne sont plus vérifiées, alors compteurBreakLoop==0 et cassera la boucle.
	//Autrement dit, Les deux couleur sont selectionnées définitivement des deux côtés.
	for {

		compteurBreakLoop = 0
		if len(selectCircle1) <= 1 {
			selectCircle1 = readFromClient(conn, in)
			writeForClient(conn2, out2, string(selectCircle1))
			compteurBreakLoop++

		}
		if len(selectCircle2) <= 1 {
			selectCircle2 = readFromClient(conn2, in2)
			writeForClient(conn, out, string(selectCircle2))
			compteurBreakLoop++
		}
		if compteurBreakLoop == 0 {
			break
		}
	}

	//Le serveur va choisir au hasard, qui sera le joueur 1 et le joueur 2
	//il va évidemment l'indiquer aux joueurs.

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(2)
	

	//Permet de savoir qui commence une partie.
	whoPlayFirst := ""

	//Si jamais C'est le client 1 qui commence
	if randomNumber == 0 {
		writeForClient(conn, out, "UR_PLAYER1")
		writeForClient(conn2, out2, "UR_PLAYER2")

		whoPlayFirst = "client1"
		//Si jamais c'est le client 2 qui commence.
	} else {
		writeForClient(conn, out, "UR_PLAYER2")
		writeForClient(conn2, out2, "UR_PLAYER1")

		whoPlayFirst = "client2"
	}





	//Boucle while pour relancer la partie
	for {
		//Si jamais le jeu se relance, il faudra voir si le client1 perd, si c'est le cas il passe ici. (car le perdant commence)
		//Si il s'agit de la première partie, alors c'est client1 qui commence si randomNumber==0.
		if whoPlayFirst == "client1" {
			//Boucle while pour de jeu
			for {
				//Lecture de client1
				placeduJetonJoueur1 := readFromClient(conn, in)

				//On vérifie si l'info qui est lue du joueur montre une fin de partie (c'est à dire qu'il perde). Pourquoi?
				//Le joueur qui gagne mets son dernier jeton, les coordonnées sont lues par le serveur. Puis elles sont transmises à l'adversaire.
				//il renverra son resultat de partie (END_GAME_I_LOSE forcément) toujours en premier.

				//Cas Exceptionnel :
				//Si le client 1 commence, lors d'une égalité, étant donné un nombre pair de case (7*6 = 42)
				//Le dernier joueur à jouer sera client 2. En suivant la logique expliquée au dessu, ça sera "END_GAME_EQUALITY" du client 1 qui sera lu en premier.

				if placeduJetonJoueur1 == "END_GAME_I_LOSE" || placeduJetonJoueur1 == "END_GAME_EQUALITY" {

					//On récupère l'info du joueur adverse(END_GAME_I_WIN). On ne l'utilisera pas, mais on le fait lire au serveur, pour avoir une symétrie des informations.
					readFromClient(conn2, in2)
					whoPlayFirst = "client1"

					//La boucle de jeu s'arrête.
					break
				}

				//Envoie à client2 de la position du jeton de client1
				writeForClient(conn2, out2, placeduJetonJoueur1)
				//Lecture de client2
				placeduJetonJoueur2 := readFromClient(conn2, in2)

				//Idem que au-dessu
				if  placeduJetonJoueur2 == "END_GAME_I_LOSE"  {
					readFromClient(conn, in)
					whoPlayFirst = "client2"

					break
				}

				//Puis écriture au client 1 et ainsi de suite...
				writeForClient(conn, out, placeduJetonJoueur2)
			}
			//Protocole dans lequel le premier client ayant dit qu'il veut rejouer dis à l'autre client qu'il est prêt.
			//Ce protocole est identique à celui de l'envoi "OK_COLOR" plus au dessu dans le main().
			//Si il accepte, la partie se relance.
			rejoue := readFromClient(conn, in)
			rejoue2 := readFromClient(conn2, in2)
			writeForClient(conn, out, string(rejoue2))
			writeForClient(conn2, out2, string(rejoue))
			continue
		}

		//Si jamais le jeu se relance, il faudra voir si le joueur2 perd, si c'est le cas il passe ici.
		if whoPlayFirst == "client2" {

			//Boucle pour jouer
			//la boucle est similaire, on commence juste à jouer par le joueur 2 au lieu du joueur 1.
			//Si le client 1 commence : for {Read client1 -> Write client 2 -> read Client 2 -> Write client 1}
			//Si le client 2 commence : for {Read client2 -> Write client 1 -> read Client 1 -> Write client 2}
			for {
				placeduJetonJoueur2 := readFromClient(conn2, in2)
				if placeduJetonJoueur2 == "END_GAME_I_LOSE" || placeduJetonJoueur2 == "END_GAME_EQUALITY"{
					readFromClient(conn, in)
					whoPlayFirst = "client2"

					break
				}
				writeForClient(conn, out, placeduJetonJoueur2)
				placeduJetonJoueur1 := readFromClient(conn, in)
				if placeduJetonJoueur1 == "END_GAME_I_LOSE"{
					readFromClient(conn2, in2)
					whoPlayFirst = "client1"

					break
				}
				writeForClient(conn2, out2, placeduJetonJoueur1)
			}

			rejoue := readFromClient(conn, in)
			rejoue2 := readFromClient(conn2, in2)
			writeForClient(conn, out, string(rejoue2))
			writeForClient(conn2, out2, string(rejoue))
			continue
		}

	}

}




// FONCTION UTILITAIRES

//Ecris au client.
func writeForClient(c net.Conn, out *bufio.Writer, message string) {

	_, err := out.WriteString(message + "\n")
	if err != nil {
		log.Println("WriteString error", err)
	}
	err = out.Flush()
	if err != nil {
		log.Println("Flush error:", err)
	}

	log.Println("Envoyé:" + message)
	return

}

//Lis depuis le client.
func readFromClient(c net.Conn, in *bufio.Reader) (message string) {

	recv, err := in.ReadString(byte('\n'))

	if err != nil {
		log.Println("ReadString error:", err)
	}

	log.Println("Message reçu:", recv)
	//Permet d'enlever le \n à la fin de la string reçue.
	recv = strings.Replace(recv, "\n", "", -1)

	return recv
}
