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
func StorageNodeListen(w http.ResponseWriter, r *http.Request){

	// Get node ID
	var nodeID int =int(r.Host[len(r.Host)-1]-'0')

	// Listen to tracker
	if r.Method == http.MethodPost{

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
		fmt.Println(chunk.Hash)

		httpVar.MapMutex.Lock()

		// if data directory doesn't exist, create it
		_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data")
		if err != nil {
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data",0777)
		}

		// if data/chunk.Hash directory doesn't exist, create it
		_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+chunk.Hash)
		if err != nil {
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+chunk.Hash,0777)
		}

		// if data/chunk.Hash/nodeID directory doesn't exist, create it
		_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID))
		if err != nil {
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID),0777)
		}

		/*
		_, exists := httpVar.HashMap[chunk.Hash]
		if !exists {
			httpVar.HashMap[chunk.Hash][nodeID-1] = true
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+chunk.Hash,07777)
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID),07777)

			// --> HERE
		}
		*/
		httpVar.MapMutex.Unlock()

		// Save chunk to file
		err=ioutil.WriteFile(path+"/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID)+"/NEW"+strconv.Itoa(httpVar.CurrentPart),[]byte(chunk.Text),0777)
		if err != nil {
			fmt.Println("StorageNodeListen: error creating/writing file", err.Error())
		}

		httpVar.TrackerMutex.Lock()
		httpVar.CurrentPart++
		httpVar.TrackerMutex.Unlock()


		// Send chunk to peers
		for _, peer :=range httpVar.Peers {
			peerURL := "http://" + peer + "/p2pRequest"
			go func(p string, URL string) {
				if  nodeID == int(p[len(p)-1]-'0'){
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
					_, err := http.Post(peerURL, "application/json", rpipe)
					httpVar.SendMutex.Unlock()
					//fmt.Println(peerURL)
					if err != nil {
						fmt.Println("Error sending http POST p2p", err.Error())
					}
				}
			}(peer, peerURL)
		}
		fmt.Println(httpVar.CurrentPart)
		if httpVar.CurrentPart == (totalPartsNum*3)-1 {
			fmt.Println("..........................................Peer END ....................................................", time.Since(start))
		}
	}



}





// Listen to other peers
func p2pRequest(w http.ResponseWriter, r *http.Request) {
	// Get peer ID
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')

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
