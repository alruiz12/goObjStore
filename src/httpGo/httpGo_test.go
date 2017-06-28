package httpGo
import (
	//"net/http/httptest"
	"testing"
	"fmt"
	"net/http"
	"time"
	"os"
)
func TestDirecories(t *testing.T){
	router := MyNewRouter()





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

	//var proxy1Addr = "127.0.0.1:8070"
	var proxy2Addr = "127.0.0.1:8071"
	var proxy3Addr = "127.0.0.1:8072"
	var proxy4Addr = "127.0.0.1:8073"
	var proxy5Addr = "127.0.0.1:8074"

	var ProxyAddr=[]string{/*proxy1Addr,*/proxy2Addr,proxy3Addr,proxy4Addr,proxy5Addr}


	peer1:=&http.Server{Addr:":8081", Handler:router}
	peer2:=&http.Server{Addr:":8082", Handler:router}
	peer3:=&http.Server{Addr:":8083", Handler:router}
	/*peer4:=&http.Server{Addr:"8084", Handler:router}
	peer5:=&http.Server{Addr:"8085", Handler:router}
	*/

	proxy1:=&http.Server{Addr:":8071", Handler:router}
	proxy2:=&http.Server{Addr:":8072", Handler:router}
	proxy3:=&http.Server{Addr:":8073", Handler:router}
	/*proxy4:=&http.Server{Addr:"8074", Handler:router}
	proxy5:=&http.Server{Addr:"8075", Handler:router}
	*/


	proxy:=&http.Server{Addr:":8080", Handler:router}
	tracker:=&http.Server{Addr:":8070", Handler:router}

	go func(){
		if err := peer1.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer2.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := peer3.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := proxy1.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := proxy2.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := proxy3.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := tracker.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := proxy.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()


	StartTracker(Peers, ProxyAddr)
	go func(){
		time.Sleep(5*time.Second)
		Put(filePath, proxyAddr, trackerAddr, 3)
		time.Sleep(5*time.Second)
		Get("0527cbea2805d89c6d5d6457b7f9f77c",ProxyAddr, trackerAddr)
	}()

	time.AfterFunc(150 * time.Second, func(){
		if err:= peer1.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy1.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy2.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy3.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= tracker.Shutdown(nil); err!=nil{
			panic(err)
		}

	})
	time.Sleep(180*time.Second)

}