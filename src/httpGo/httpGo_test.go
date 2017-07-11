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
	router :=MyNewRouter()
	const IP = "127.0.0.1"


	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/dataset.xml"

	// TRACKER
	var trackerAddr = IP+":8000"

	// PEERS
	// IMPORTANT NOTE: last character of port is the peer's internal identifier
	var Peer1a = IP+":8011"
	var Peer1b = IP+":8021"
	var Peer1c = IP+":8031"
	var peer1List = []string{Peer1a, Peer1b, Peer1c}

	var Peer2a = IP+":8012"
	var Peer2b = IP+":8022"
	var Peer2c = IP+":8032"
	var peer2List = []string{Peer2a, Peer2b, Peer2c}

	var Peer3a = IP+":8013"
	var Peer3b = IP+":8023"
	var Peer3c = IP+":8033"
	var peer3List = []string{Peer3a, Peer3b, Peer3c}

	var Peers =[][]string{peer1List, peer2List, peer3List/*, peer4List, peer5List*/}

	// PROXY
	var proxy1 = IP+":8070"
	var proxy2 = IP+":8071"
	var proxy3 = IP+":8072"

	var ProxyAddr=[]string{proxy1,proxy2,proxy3/*,proxy4,proxy5*/}






	peer1arun:=&http.Server{Addr:Peer1a, Handler:router}
	peer1brun:=&http.Server{Addr:Peer1b, Handler:router}
	peer1crun:=&http.Server{Addr:Peer1c, Handler:router}
	peer2arun:=&http.Server{Addr:Peer2a, Handler:router}
	peer2brun:=&http.Server{Addr:Peer2b, Handler:router}
	peer2crun:=&http.Server{Addr:Peer2c, Handler:router}
	peer3arun:=&http.Server{Addr:Peer3a, Handler:router}
	peer3brun:=&http.Server{Addr:Peer3b, Handler:router}
	peer3crun:=&http.Server{Addr:Peer3c, Handler:router}
	/*peer4:=&http.Server{Addr:"8084", Handler:router}
	peer5:=&http.Server{Addr:"8085", Handler:router}
	*/

	proxy1run:=&http.Server{Addr:proxy1, Handler:router}
	proxy2run:=&http.Server{Addr:proxy2, Handler:router}
	proxy3run:=&http.Server{Addr:proxy3, Handler:router}
	/*proxy4:=&http.Server{Addr:"8074", Handler:router}
	proxy5:=&http.Server{Addr:"8075", Handler:router}
	*/



	tracker:=&http.Server{Addr:trackerAddr, Handler:router}

	go func(){
		if err := peer1arun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer1brun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer1crun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()

	go func(){
		if err := peer2arun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer2brun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer2crun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()

	go func(){
		if err := peer3arun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer3brun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer3crun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()

	go func(){
		fmt.Println("go func")
		if err := proxy1run.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := proxy2run.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := proxy3run.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		fmt.Println("go func")
		if err := tracker.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()


	StartTracker(Peers)

	go func(){
		Put(filePath, trackerAddr, 3)

		time.Sleep(5*time.Second)

		Get("0527cbea2805d89c6d5d6457b7f9f77c",ProxyAddr, trackerAddr)

		time.Sleep(45*time.Second)

		if (CheckPieces("0527cbea2805d89c6d5d6457b7f9f77c","NEW.xml",filePath, 3)) == false {
			t.Error("True expected")
		}
	}()

	time.AfterFunc(600 * time.Second, func(){
		if err:= peer1arun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2arun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3arun.Shutdown(nil); err!=nil{
			panic(err)
		}

		if err:= peer1brun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2brun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3brun.Shutdown(nil); err!=nil{
			panic(err)
		}

		if err:= peer1crun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2brun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3crun.Shutdown(nil); err!=nil{
			panic(err)
		}


		if err:= proxy1run.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy2run.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy3run.Shutdown(nil); err!=nil{
			panic(err)
		}

		if err:= tracker.Shutdown(nil); err!=nil{
			panic(err)
		}

	})
	time.Sleep(180*time.Second)

}
