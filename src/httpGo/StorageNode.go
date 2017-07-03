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
	"path/filepath"
	"strings"
	"sync"

)
var path = (os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT")
var chunk msg
var hashRecived hashMsg
func prepNode(w http.ResponseWriter, r *http.Request){
	var nodeID int =int(r.Host[len(r.Host)-1]-'0')
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ",err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ",err)
	}
	if err := json.Unmarshal(body, &hashRecived); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ",err)
		}
	}

	httpVar.DirMutex.Lock()
	// if data directory doesn't exist, create it
	_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data")
	if err != nil {
		os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data",0777)
	}

	// if data/chunk.Hash directory doesn't exist, create it
	_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+hashRecived.Hash)
	if err != nil {
		os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+hashRecived.Hash,0777)
	}

	// if data/chunk.Hash/nodeID directory doesn't exist, create it
	_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+hashRecived.Hash+"/"+strconv.Itoa( nodeID))
	if err != nil {
		err2:=os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+hashRecived.Hash+"/"+strconv.Itoa( nodeID),0777)
		if err2!=nil{
			fmt.Println("StorageNode error making dir", err.Error())
		}
	}
	httpVar.DirMutex.Unlock()
}

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

		 var wg sync.WaitGroup
                wg.Add(len(chunk.NodeList)+1)

		go func(CurrentPart int, text string){
			// Save chunk to file
			err=ioutil.WriteFile(path+"/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID)+"/NEW"+strconv.Itoa(CurrentPart),[]byte(text),0777)
			if err != nil {
				fmt.Println("StorageNodeListen: error creating/writing file", err.Error())
			}
			defer wg.Done()
		}(chunk.Name, chunk.Text)

		httpVar.TrackerMutex.Lock()
		httpVar.CurrentPart++
		httpVar.TrackerMutex.Unlock()
		var currentAdr int = 0
		// Send chunk to peers
		for _, peer :=range chunk.NodeList {
			peerURL := "http://" + peer[currentAdr] + "/p2pRequest"

			go func(p []string, URL string) {
				if  nodeID == int(p[0][len(p)-1]-'0'){	// Don't send to itself

				}else {
					rpipe, wpipe :=io.Pipe() // create pipe
					go func(){
						err:=json.NewEncoder(wpipe).Encode(&chunk)
						wpipe.Close()			// close pipe when go routine finishes
						if err != nil {
							fmt.Println("Error encoding to pipe ", err.Error())
						}
					}()
					//httpVar.SendMutex.Lock()
				 //	httpVar.SendReady <- 1
					_, err := http.Post(peerURL, "application/json", rpipe)
				//	<- httpVar.SendReady
					//httpVar.SendMutex.Unlock()
					if err != nil {
						fmt.Println("Error sending http POST p2p", err.Error())
					}
				}

				defer wg.Done()
			}(peer, peerURL)
			currentAdr=(currentAdr+1)%len(chunk.NodeList[0])
		}
		wg.Wait()
		if httpVar.CurrentPart == (totalPartsNum*chunk.Num)-1 {
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
		httpVar.DirMutex.Lock()

		err = ioutil.WriteFile(path + "/src/data/"+chunk.Hash+"/"+strconv.Itoa( peerID)+ "/P2P" + strconv.Itoa(httpVar.P2pPart), []byte(chunk.Text), 0777)
		if err != nil {
			fmt.Println("P2pRequest: Peer: error creating/writing file p2p", err.Error())
		}
		_, err = os.Stat(path+ "/src/data/"+chunk.Hash+"/"+strconv.Itoa( peerID)+ "/P2P" + strconv.Itoa(httpVar.P2pPart))
		if err != nil {
			err = ioutil.WriteFile(path + "/src/data/"+chunk.Hash+"/"+strconv.Itoa( peerID)+ "/P2P" + strconv.Itoa(httpVar.P2pPart), []byte(chunk.Text), 0777)
                if err != nil {
                        fmt.Println("P2pRequest: Peer: error creating/writing file p2p", err.Error())
                }
		}
		httpVar.DirMutex.Unlock()
		if httpVar.P2pPart >= (totalPartsNum*(chunk.Num-1))-1 {
			fmt.Println("p2p: ",time.Since(start))
			fmt.Println("..........................................p2p END ....................................................",httpVar.P2pPart)
			return
		}

		httpVar.PeerMutex.Lock()
		httpVar.P2pPart++
		httpVar.PeerMutex.Unlock()

	}
}

func GetChunks(w http.ResponseWriter, r *http.Request){
	// Get node ID
	var nodeIDint int = int(r.Host[len(r.Host) - 1] - '0')
	var nodeID string = strconv.Itoa(nodeIDint)
	var keyURL jsonKeyURL
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &keyURL); err != nil {
		fmt.Println("GetChunks error Unmashalling ",err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(423) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("GetChunks: error unprocessable entity: ",err.Error())
			return
		}
		return
	}


	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	go sendChunksToProxy(nodeID, keyURL.Key, keyURL.URL)
}

type getMsg struct {
	Text string
	Name string
	NodeID string
	Key string
}
func sendChunksToProxy(nodeID string, key string, URL string){
	//partBuffer:=make([]byte,fileChunk)
	proxyURL:="http://"+URL

	// for each proxy-name in directory send
	filepath.Walk(path + "/src/data/"+key+"/"+nodeID, func(path string, info os.FileInfo, err error) error {

		if strings.Contains(info.Name(),"NEW") {
			partBuffer:=make([]byte,info.Size())
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("sendChunksToProxy error opening file ",err.Error())
			}
			_, err = file.Read(partBuffer)
			if err != nil {
				fmt.Println("sendChunksToProxy error opening file ",err.Error())
			}
			m:=getMsg{Text: string(partBuffer), Name: info.Name(), NodeID:nodeID, Key:key}
			r, w :=io.Pipe()			// create pipe
			go func() {
				defer w.Close()			// close pipe when go routine finishes
				// save buffer to object
				err=json.NewEncoder(w).Encode(&m)
				if err != nil {
					fmt.Println("Error encoding to pipe ", err.Error())
				}
			}()
			res, err := http.Post(proxyURL,"application/json", r )
			if err != nil {
				fmt.Println("sendChunksToProxy: error creating request: ",err.Error())
			}
			if err := res.Body.Close(); err != nil {
				fmt.Println(err)
			}
		}

		return nil

	})
}
