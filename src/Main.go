package main

import (
	"github.com/alruiz12/simpleBT/src/tcp"
	"os"
	"time"
)

func main() {
	// tracker sends to:
	var peer1 = 8051
	var peer2 = 8052
	var peer3 = 8053
	peers :=[]int{peer1, peer2, peer3}

	// peers p2p address, must start with ":809"
	var peer1p2p = 8091
	var peer2p2p = 8092
	var peer3p2p = 8093

	peersp2p :=[]int{peer1p2p, peer2p2p, peer3p2p}

	// last character of port is the peer's internal identifier
	// e.g. peer1 has port ":8081" and port ":8091"

	go func() {
		tcp.PeerListen(127,0,0,1,peer1,peersp2p)
	}()
	go func() {
		tcp.PeerListen(127,0,0,1,peer2,peersp2p)
	}()
	go func() {
		tcp.PeerListen(127,0,0,1,peer3,peersp2p)
	}()
	time.Sleep(2 * time.Second)
	//tcp.TrackerFile("127.0.0.1",peers, os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")
	//tcp.TrackerFile("127.0.0.1",peers,os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/dataset.xml")
	tcp.TrackerDivideLoad(127,0,0,1,peers, os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")


	}
