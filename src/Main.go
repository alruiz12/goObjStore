package main

import (
	"os"
	//"time"
	"github.com/alruiz12/simpleBT/src/httpGo"
	"net/http"
	"github.com/alruiz12/simpleBT/src/httpVar"
)

func main() {
	//var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile"
	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/dataset.xml"


	var trackerAddr = ":8080"

	// tracker sends to:
	httpVar.Peer1 = "127.0.0.1:8081"
	httpVar.Peer2 = "127.0.0.1:8082"
	httpVar.Peer3 = "127.0.0.1:8083"
	httpVar.Peers =[]string{httpVar.Peer1, httpVar.Peer2, httpVar.Peer3}
	// last character of port is the peer's internal identifier

	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()
	go func(){httpGo.TrackerDivideLoad(filePath, trackerAddr ,httpVar.Peers)}()
	go func(){http.ListenAndServe(":8081", routerPeer)}()
	go func(){http.ListenAndServe(":8082", routerPeer)}()
	go func(){http.ListenAndServe(":8083", routerPeer)}()
	http.ListenAndServe(":8080", routerTracker)
	}

