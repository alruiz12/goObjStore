package httpGo
import(
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"os"
	"encoding/json"
	"github.com/alruiz12/simpleBT/src/httpVar"
	"strconv"
	"time"
)
var path = (os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT")
var chunk msg
func PeerListen(w http.ResponseWriter, r *http.Request){
	// Get peer ID
	var peerID int =int(r.Host[len(r.Host)-1]-'0')
	// Create folder for chunks


	// Listen to tracker
	if r.Method == http.MethodPost{
		os.Mkdir(path+"/src/httpReceived"+strconv.Itoa(peerID),07777)

		fmt.Println( "gettig post" )

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ",err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ",err)
		}
		if err := json.Unmarshal(body, &chunk); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			log.Println(err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ",err)
			}
		}

		// Save chunk to file
		err=ioutil.WriteFile(path+"/src/httpReceived"+strconv.Itoa(peerID)+"/NEW"+strconv.Itoa(httpVar.CurrentPart),[]byte(chunk.Text),0777)
		if err != nil {
			fmt.Println("Peer: error creating/writing file", err.Error())
		}

		// Send chunk to peers

		httpVar.TrackerMutex.Lock()
		httpVar.CurrentPart++
		httpVar.TrackerMutex.Unlock()


	}





	for _, peer :=range httpVar.Peers {
		peerURL := "http://" + peer + "/p2pRequest"
		go func(p string, URL string	) {
			if peerID == int(p[len(p)-1]-'0'){
			}else {
				rpipe, wpipe :=io.Pipe() // create pipe
				go func(){
				err:=json.NewEncoder(wpipe).Encode(&chunk)
				wpipe.Close()			// close pipe when go routine finishes
				if err != nil {
					fmt.Println("Error encoding to pipe ", err.Error())
				}
				}()
				httpVar.SendMutex.Lock()
				//peerURL2 := "http://" + peer + "/p2pRequest"


				_, err := http.Post(peerURL, "application/json", rpipe)
				httpVar.SendMutex.Unlock()
				fmt.Println(peerURL)
				if err != nil {
					fmt.Println("Error sending http POST p2p", err.Error())
				}
			}
		}(peer, peerURL)
		fmt.Println("-------------------------------------------------------------------------------------",peer)
	}

	fmt.Println("*** addFile FINISHES ***")

	if httpVar.CurrentPart == totalPartsNum-1 {
		fmt.Println(time.Since(start))
		fmt.Println("..........................................Peer END ....................................................")
	}


	// Listen to other peers
}






func p2pRequest(w http.ResponseWriter, r *http.Request) {
	// Get peer ID
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')
	// Create folder for chunks
	fmt.Println("p2p!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	// Listen to tracker
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ", err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ", err)
		}
		if err := json.Unmarshal(body, &chunk); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			log.Println(err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ", err)
			}
		}

		// Save chunk to file
		err = ioutil.WriteFile(path + "/src/httpReceived" + strconv.Itoa(peerID) + "/P2P" + strconv.Itoa(httpVar.P2pPart), []byte(chunk.Text), 0777)
		if err != nil {
			fmt.Println("Peer: error creating/writing file p2p", err.Error())
		}

		if httpVar.P2pPart >= (totalPartsNum*6)-1 {
			fmt.Println("p2p: ",time.Since(start))
			fmt.Println("..........................................p2p END ....................................................",httpVar.P2pPart)
			return
		}

		httpVar.PeerMutex.Lock()
		httpVar.P2pPart++
		httpVar.PeerMutex.Unlock()

	}
}
