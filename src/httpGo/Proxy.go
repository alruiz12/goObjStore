package httpGo
import(
	"os"
	"fmt"
	"strconv"
	"math"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
	"mime/multipart"
	"time"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/alruiz12/simpleBT/src/httpVar"
)

//const fileChunk = 1*(1<<10) // 1 KB
const fileChunk = 8*(1<<20) // 8 MB

type msg struct {
	NodeList []string
	Num int
	Hash string
	Text string
}
var totalPartsNum int
var start time.Time
func ProxyDivideLoad(filePath string, addr string, trackerAddr string, numNodes int){
	time.Sleep(1 * time.Second)
	start=time.Now()
	var err error

	// ask tracker for nodes
	quantityJson := `{"Quantity":"`+strconv.Itoa(numNodes)+`"}`
	reader := strings.NewReader(quantityJson)
	trackerURL:="http://"+trackerAddr+"/GetNodes"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("ProxyDivideLoad: error creating request: ",err.Error())
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("ProxyDivideLoad: error sending request: ",err.Error())
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := res.Body.Close(); err != nil {
		panic(err)
	}
	var nodeList []string
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("ProxyDivideLoad: error unprocessable entity: ",err.Error())
		return
	}
	fmt.Println("nodeList:",nodeList[0])
	if err != nil {
		fmt.Println("ProxyDivideLoad: error reciving response: ",err.Error())
	}



	var currentPart int = 0
	//var currentPeer string
	var partSize int
	var currentNum int = 0
	var partBuffer []byte
	var peerURL string
	//var body *bytes.Buffer
	var writer *multipart.Writer
	var buf bytes.Buffer
	_,_=writer, buf // avoiding declared but not used
	var hash string = md5sum(filePath)
	var auxList []bool
	var i int = 0
	fmt.Println(numNodes)
	for i<numNodes {
		fmt.Println(numNodes)
		auxList=append(auxList, false)
		i++
	}
	httpVar.DirMutex.Lock()
	httpVar.HashMap[hash]=auxList
	httpVar.DirMutex.Unlock()
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
		//fmt.Println("					CURRENT ",currentPart)
		//currentPeer=peers[currentNum]
		partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
		partBuffer=make([]byte,partSize)
		_,err = file.Read(partBuffer)		// Get chunk
		m:=msg{NodeList:nodeList, Num:numNodes, Hash:hash, Text:string(partBuffer)}
		r, w :=io.Pipe()			// create pipe

		go func() {
			defer w.Close()			// close pipe when go routine finishes
			 // save buffer to object
			err=json.NewEncoder(w).Encode(&m)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
			}
		}()


		peerURL = "http://"+nodeList[currentNum]+"/StorageNodeListen"
		// Send to current peer
		_, err := http.Post(peerURL, "application/json", r )
		if err != nil {
			fmt.Println("Error sending http POST ", err.Error())
		}

		currentPart++
		fmt.Println("adsfasfasfdsdaf ", chunk.Num)
		currentNum=(currentNum+1)%chunk.Num
	}
	fmt.Println("..........................................Proxy END ....................................................")
	}


/*
md5sum opens the file we want to compute the hash and computes it
@param path to the file we want to split
returns the computed hash
*/
func md5sum(filePath string) string{
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash,file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	mainFileHash:=hex.EncodeToString(hash.Sum(nil))
	return mainFileHash
}