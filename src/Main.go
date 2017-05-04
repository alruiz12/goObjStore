package main

import (
	/*
	"github.com/alruiz12/simpleBT/src/conf"
	//"os"
	//"fmt"

	"net/http"
	"log"
	"github.com/alruiz12/simpleBT/src/vars"
	"time"
	"net"
	"strings"


	"log"
	"os"
	"os/exec"
	*/
	"github.com/alruiz12/simpleBT/src/tcp"
	"time"
	"os"
)

func main() {

	//vars.IP = ""
	//conf.SplitFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")
	//fmt.Println(conf.CheckPieces("bigFile"))

	/*
	router := conf.MyNewRouter()
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		if iface.Flags & net.FlagUp == 0 {
			continue //interface down
		}
		if iface.Flags & net.FlagLoopback != 0 {
			continue //loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		if strings.Compare(iface.Name, "enp0s8") == 0 {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type){
				case *net.IPNet: ip = v.IP
				case *net.IPAddr: ip = v.IP
				}
				ip = ip.To4()
				vars.IP = ip.String()
				break

			}
			break
		}

	}

	vars.TrackerIP="10.0.0.10"
	if strings.Compare(vars.IP, "10.0.0.10") != 0 {
		go func() {
			var quit = make(chan int)
			conf.StartAnnouncing(2, 20, "10.0.0.12", "bigFile", quit)
			time.AfterFunc(90 * time.Second, func() {
				close(quit)
			})
		}()
	} else {

		go func() {
			conf.CheckInactivePeers(5)
		}()
	}
	log.Fatal(http.ListenAndServe(vars.TrackerPort, router))

	*/
	go func() {
		tcp.PeerListen("127.0.0.1:8081")
	}()
	time.Sleep(2 * time.Second)
	tcp.TrackerFile("127.0.0.1",":8081", os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/bigFile")

}
