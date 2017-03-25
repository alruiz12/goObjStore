package main

import (
	"net/http"
	"log"
	"github.com/alruiz12/simpleBT/src/conf"
	"github.com/alruiz12/simpleBT/src/vars"
	"time"
	"net"
	"strings"
)

func main() {
	var IP string
	IP=""
	router := conf.MyNewRouter()
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces{
		if iface.Flags&net.FlagUp ==0{
			continue //interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue //loopback interface
		}
		addrs, err:= iface.Addrs()
		if err != nil {
			panic(err)
		}
		if strings.Compare( iface.Name, "enp0s8")==0 {
			for _, addr := range addrs{
				var ip net.IP
				switch v := addr.(type){
					case *net.IPNet: ip=v.IP
					case *net.IPAddr: ip=v.IP
				}
				ip=ip.To4()
				IP=ip.String()
				break

			}
			break
		}

	}
	go func() {
		var quit = make(chan int)
		conf.StartAnnouncing(2,9,IP,"torrent1",quit)
		time.AfterFunc(9 * time.Second, func(){close(quit)})
	}()
	/*
	// Seeder
	go func() {
		var quit = make(chan int)
		conf.StartAnnouncing(2,9,"192.168.0.1","torrent1",quit)
		time.AfterFunc(9 * time.Second, func(){close(quit)})
	}()

	// Peer1
	go func() {
		var quit = make(chan int)
		conf.StartAnnouncing(2,15,"192.168.0.2","torrent1", quit)
		time.AfterFunc(15 * time.Second, func(){close(quit)})

	}()

	// Peer2
	go func(){
		var quit = make(chan int)
		conf.StartAnnouncing(2,21,"192.168.0.3","torrent1",quit )
		time.AfterFunc(21 * time.Second, func(){close(quit)})
	}()
	*/
	go func() {
		conf.CheckInactivePeers(5)
	}()
	log.Fatal(http.ListenAndServe(vars.TrackerPort, router))
}
