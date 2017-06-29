package main

import (
	"os"
	"time"
	"github.com/alruiz12/simpleBT/src/httpGo"
	"net/http"
	"fmt"

)

func main() {


	//var filePath2 = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile"
	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/dataset.xml"


	var proxyAddr = "127.0.0.1:8080"
	var trackerAddr = "127.0.0.1:8070"

	// tracker sends to:
	var Peer1 = "127.0.0.1:8081"
	var Peer2 = "127.0.0.1:8082"
	var Peer3 = "127.0.0.1:8083"
	var Peer4 = "127.0.0.1:8084"
	var Peer5 = "127.0.0.1:8085"
	var Peers =[]string{Peer1, Peer2, Peer3, Peer4, Peer5}
	// last character of port is the peer's internal identifier

	//var proxy1 = "127.0.0.1:8070"
	var proxy2 = "127.0.0.1:8071"
	var proxy3 = "127.0.0.1:8072"
	var proxy4 = "127.0.0.1:8073"
	var proxy5 = "127.0.0.1:8074"

	var ProxyAddr=[]string{/*proxy1,*/proxy2,proxy3,proxy4,proxy5}

	httpGo.StartTracker(Peers, ProxyAddr)



	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()
	go func(){http.ListenAndServe(":8070", routerPeer)}()
	go func(){httpGo.Put(filePath, proxyAddr, trackerAddr, 3)
		time.Sleep(5*time.Second)

		httpGo.Get("0527cbea2805d89c6d5d6457b7f9f77c",ProxyAddr, trackerAddr)
		fmt.Println(httpGo.CheckPieces("0527cbea2805d89c6d5d6457b7f9f77c","NEW"))

	}()
	go func(){http.ListenAndServe(":8081", routerPeer)}()
	go func(){http.ListenAndServe(":8082", routerPeer)}()
	go func(){http.ListenAndServe(":8083", routerPeer)}()
	go func(){http.ListenAndServe(":8084", routerPeer)}()
	go func(){http.ListenAndServe(":8085", routerPeer)}()


	go func(){http.ListenAndServe(":8070", routerPeer)}()
	go func(){http.ListenAndServe(":8071", routerPeer)}()
	go func(){http.ListenAndServe(":8072", routerPeer)}()
	go func(){http.ListenAndServe(":8073", routerPeer)}()
	go func(){http.ListenAndServe(":8074", routerPeer)}()


	http.ListenAndServe(":8080", routerTracker)

	}

