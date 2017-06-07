package main

import (
	"os"
	//"time"
	"github.com/alruiz12/simpleBT/src/httpGo"
	"net/http"

)

func main() {



	//var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile"
	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/dataset.xml"
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/httpReceived1",07777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/httpReceived2",07777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/httpReceived3",07777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data",07777)

	var proxyAddr = "127.0.0.1:8080"
	var trackerAddr = "127.0.0.1:8070"

	// tracker sends to:
	var Peer1 = "127.0.0.1:8081"
	var Peer2 = "127.0.0.1:8082"
	var Peer3 = "127.0.0.1:8083"
	var Peers =[]string{Peer1, Peer2, Peer3}
	// last character of port is the peer's internal identifier


	httpGo.StartTracker(Peers)


	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()
	go func(){http.ListenAndServe(":8070", routerPeer)}()
	go func(){httpGo.ProxyDivideLoad(filePath, proxyAddr, trackerAddr, 3)}()
	go func(){http.ListenAndServe(":8081", routerPeer)}()
	go func(){http.ListenAndServe(":8082", routerPeer)}()
	go func(){http.ListenAndServe(":8083", routerPeer)}()
	http.ListenAndServe(":8080", routerTracker)

	}

