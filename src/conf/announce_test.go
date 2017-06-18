package conf

import (
	"testing"

	"net/http"
	"github.com/alruiz12/simpleBT/src/vars"
	"log"
	"time"
	"fmt"
)

func TestAnnounce(t *testing.T){
	router := MyNewRouter()

	srv:=&http.Server{Addr: vars.TrackerPort, Handler:router}
	go func(){
		fmt.Println("go func")
		if err := srv.ListenAndServe(); err!=nil{
			log.Printf("ListenAndServe error", err)
		}
	}()
	time.Sleep(3*time.Second)
	go func() {
		var quit = make(chan int)
		StartAnnouncing(2, 11, vars.IP, "torrent1", quit)
		time.AfterFunc(9 * time.Second, func() {
			close(quit)
		})
	}()

	go func() {
		CheckInactivePeers(5)
	}()

	time.AfterFunc(15 * time.Second, func(){
		if err:= srv.Shutdown(nil); err!=nil{
			panic(err)
		}
	})
	time.Sleep(18*time.Second)




}
