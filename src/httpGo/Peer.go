package httpGo
import(
	"net/http"
	"fmt"
	"io/ioutil"
	//"io"
	"log"
	"os"
	"encoding/json"
	"github.com/alruiz12/simpleBT/src/httpVar"
	"strconv"
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
		httpVar.Mutex.Lock()
		httpVar.CurrentPart++
		httpVar.Mutex.Unlock()


	}


	fmt.Println("*** addFile FINISHES ***")


		// Send chunk to peers

	// Listen to other peers
}
