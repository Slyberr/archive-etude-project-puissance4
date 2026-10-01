package main

import (
	"log"
	"net"
	"bufio"
	"strings"
)


//FONCTIONS UTILITAIRES POUR LES CLIENTS

//Ecris au serveur
func writeForServer(c net.Conn,out *bufio.Writer,message string) {

	_,err:= out.WriteString(message+"\n")
	if err != nil{
		log.Println("WriteString error",err)
	}
	err=out.Flush()
	if err!=nil{
		log.Println("Flush error:",err)
	}
	log.Println("Envoyé:",message)
	return

}

//Lis depuis le serveur
func readFromServer(c net.Conn,in *bufio.Reader) (message string){

	recv,err :=in.ReadString(byte('\n'))
	if err !=nil{
		log.Println("ReadString error:",err)
	}

	log.Println("Message reçu:",recv)
	//Permet d'enlever le \n à la fin de la string reçue.
	recv=strings.Replace(recv, "\n", "", -1)

	return recv
}



//Lis depuis le serveur, cette fonction est systématiquement appelée en tant que Go Routine.
//Ainsi, elle est appellée lorsqu'un client attend l'autre sans mettre à mal le bon fonctionnement du jeu.
func readFromServerForRoutine(c net.Conn,in *bufio.Reader,channel chan string){

	//On place un verrou sur toute la fonction pour s'assurer de son intégrité.
	//Expérience: Durant plusieurs parties, nous avons parfois remarqué des pertes d'informations [Exemple Bonjour était reçu en "onj"], c'est cela qui nous a encouragé à mettre en place un verrou.
	monVerrou.Lock()
	recv,err :=in.ReadString(byte('\n'))
	if err !=nil{
		log.Println("ReadString error:",err)
	}


	log.Println("Message reçu:",recv)
	//Permet d'enlever le \n à la fin de la string reçue.
	recv=strings.Replace(recv, "\n", "", -1)
	monVerrou.Unlock()

	channel<-recv

}
