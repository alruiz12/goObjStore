package httpGo
import(
	"os"
	"fmt"
	"strconv"
	"math"
	"net/http"
	"io"
	//"io/ioutil"
	"bytes"
	"mime/multipart"
	"time"
	"encoding/json"
)

//const fileChunk = 1*(1<<10) // 1 KB
const fileChunk = 8*(1<<20) // 8 MB

type msg struct {
	Text string
}
var totalPartsNum int
var start time.Time
func TrackerDivideLoad(filePath string, addr string, peers []string){
	time.Sleep(1 * time.Second)
	start=time.Now()
	var currentPart int = 0
	var currentPeer string
	var partSize int
	var currentNum int = 0
	var partBuffer []byte
	var err error
	var peerURL string
//	var body *bytes.Buffer
	var writer *multipart.Writer
	var buf bytes.Buffer
	_,_=writer, buf // avoiding declared but not used

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	text:=strconv.FormatInt(fileInfo.Size(),10)	// size
	size,_:=strconv.Atoi(text)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	totalPartsNum= int(math.Ceil(float64(size)/float64(fileChunk)))
	fmt.Println(totalPartsNum)

	for currentPart<totalPartsNum{
		fmt.Println("					CURRRRRRRRRRRRRRRRRRRRRRR ",currentPart)
		currentPeer=peers[currentNum]
		partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
		partBuffer=make([]byte,partSize)
		_,err = file.Read(partBuffer)		// Get chunk
		r, w :=io.Pipe()			// create pipe

		go func() {
			defer w.Close()			// close pipe when go routine finishes
			m:=msg{Text:string(partBuffer)} // save buffer to object
			err=json.NewEncoder(w).Encode(&m)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
			}
		}()


		peerURL = "http://"+currentPeer+"/PeerListen"
		// Send to current peer
		_, err := http.Post(peerURL, "application/json", r )
		if err != nil {
			fmt.Println("Error sending http POST ", err.Error())
		}

		currentPart++
		currentNum=(currentNum+1)%3
	}
	fmt.Println("..........................................Tracker END ....................................................")
	fmt.Println("..........................................Tracker END ....................................................")
	fmt.Println("..........................................Tracker END ....................................................")
	fmt.Println("..........................................Tracker END ....................................................")
}