package main

import (
		"github.com/Kaibling/gogam"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)
	//gameField :=

	//gogam.ShowMap(*board)
	//log.Println((*board)[0][0])

	server := new(gogam.GameServer)
	server.StartServer()

}
