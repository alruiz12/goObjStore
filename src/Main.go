package main

import (
	"time"
	"github.com/alruiz12/simpleBT/src/httpGo"
	"net/http"
	//"fmt"
	"github.com/alruiz12/simpleBT/src/conf"

)
func main() {


	httpGo.StartTracker(conf.Peers)

	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()

	go func(){
		//httpGo.Put( conf.FilePath, conf.TrackerAddr, 3)

		time.Sleep(5*time.Second)
		//httpGo.CreateAccount("alvaro")
		//httpGo.Get(conf.KeyExample ,conf.ProxyAddr, conf.TrackerAddr)
		
		time.Sleep(45*time.Second)
		
		//fmt.Println(httpGo.CheckPieces(conf.KeyExample ,"NEW.xml",conf.FilePath, conf.NumNodes))

	}()
	go func(){http.ListenAndServe(conf.Peer1a, routerPeer)}()
	go func(){http.ListenAndServe(conf.Peer1b, routerPeer)}()
	go func(){http.ListenAndServe(conf.Peer1c, routerPeer)}()

	go func(){http.ListenAndServe(conf.Peer2a, routerPeer)}()
	go func(){http.ListenAndServe(conf.Peer2b, routerPeer)}()
	go func(){http.ListenAndServe(conf.Peer2c, routerPeer)}()

	go func(){http.ListenAndServe(conf.Peer3a, routerPeer)}()
	go func(){http.ListenAndServe(conf.Peer3b, routerPeer)}()
	go func(){http.ListenAndServe(conf.Peer3c, routerPeer)}()


	go func(){http.ListenAndServe(conf.Proxy1, routerPeer)}()
	go func(){http.ListenAndServe(conf.Proxy2, routerPeer)}()
	go func(){http.ListenAndServe(conf.Proxy3, routerPeer)}()



	http.ListenAndServe(conf.TrackerAddr, routerTracker)

	}

