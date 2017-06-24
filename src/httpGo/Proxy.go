package httpGo

import (
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
	"github.com/alruiz12/simpleBT/src/httpVar"
	"strings"
	//"sync"
)

//const fileChunk = 1*(1<<10) // 1 KB
const fileChunk = 8*(1<<20) // 8 MB

type msg struct {
	NodeList []string
	Num int
	Hash string
	Text string
	CurrentNode int
	Name int
}
var totalPartsNum int
var start time.Time
var startGet time.Time
var numGets int = 0
func Put(filePath string, addr string, trackerAddr string, numNodes int){
	time.Sleep(1 * time.Second)
	start=time.Now()
	var hash string = md5sum(filePath)
	var err error

	// ask tracker for nodes
	quantityJson := `{"Quantity":"`+strconv.Itoa(numNodes)+`","Hash":"`+hash+`"}`
	//jsonContent := `{"file":"`+torrentName+`","IP":"`+IP+`"}`
	reader := strings.NewReader(quantityJson)
	trackerURL:="http://"+trackerAddr+"/GetNodes"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Put: error creating request: ",err.Error())
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Put: error sending request: ",err.Error())
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
		fmt.Println("Put: error unprocessable entity: ",err.Error())
		return
	}
	fmt.Println("nodeList:",nodeList[0])
	if err != nil {
		fmt.Println("Put: error reciving response: ",err.Error())
	}



	var currentPart int = 0
	//var currentPeer string
	var partSize int
	var currentNum int = 0
	var partBuffer []byte
	//var peerURL string
	//var body *bytes.Buffer
	var writer *multipart.Writer
	var buf bytes.Buffer
	_,_=writer, buf // avoiding declared but not used

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
	//return
	/*var wg sync.WaitGroup
	wg.Add(totalPartsNum)*/
	for currentPart<totalPartsNum{
		partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
		partBuffer=make([]byte,partSize)
		_,err = file.Read(partBuffer)		// Get chunk
		m:=msg{NodeList:nodeList, Num:numNodes, Hash:hash, Text:string(partBuffer), CurrentNode:currentNum, Name: currentPart}
		r, w :=io.Pipe()			// create pipe
		//go func(url string, r2 io.Reader, pb string, cp int) {
			go func() {
				defer w.Close()			// close pipe when go routine finishes
				 // save buffer to object
				err=json.NewEncoder(w).Encode(&m)
				if err != nil {
					fmt.Println("Error encoding to pipe ", err.Error())
				}
			}()

			// peerURL = "http://"+nodeList[currentNum]+"/StorageNodeListen"
			// Send to current peer
			_, err := http.Post("http://"+nodeList[currentNum] + "/StorageNodeListen", "application/json", r )
			if err != nil {
				fmt.Println("Error sending http POST ", err.Error())
			}
			//fmt.Println(url, ", ", pb, "         ,", strconv.Itoa(cp ))
			//defer wg.Done()
		//}("http://"+nodeList[currentNum] + "/StorageNodeListen", r, string(partBuffer[:25]), currentPart )
		currentPart++
		currentNum=(currentNum+1)%numNodes
	}
	fmt.Println("..........................................Proxy END ....................................................")
	//wg.Wait()
	time.Sleep( 10 * time.Second)
	fmt.Println("WaitGroup waited!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
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

type jsonKeyURL struct {
	Key string		`json:"Key"`
	URL string		`json:"URL"`
}


func Get(Key string, proxyAddr []string, trackerAddr string){
	time.Sleep(1 * time.Second)

	// Ask tracker for nodes
	startGet=time.Now()
	var err error
	// ask tracker for nodes for a given key
	keyJson := `{"Key":"`+Key+`"}`
	reader := strings.NewReader(keyJson)
	trackerURL:="http://"+trackerAddr+"/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Get: error creating request: ",err.Error())
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Get: error sending request: ",err.Error())
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
		fmt.Println("Get: error unprocessable entity: ",err.Error())
		return
	}
	fmt.Println("nodeList for key: ",nodeList)
	if err != nil {
		fmt.Println("Get: error reciving response: ",err.Error())
	}
	// Create folder for receiving
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/local",+0777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/local/"+Key,0777)


	//var k jsonKey
	// For each node ask for all their Proxy-pieces
	for index, node := range nodeList {
		r, w :=io.Pipe()			// create pipe
		k:=jsonKeyURL{Key:Key, URL:proxyAddr[index]+"/ReturnData"}

		go func() {
			defer w.Close()			// close pipe when go routine finishes
			// save buffer to object
			err=json.NewEncoder(w).Encode(&k)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
			}
		}()
		url:="http://"+node+"/GetChunks"
		res, err := http.Post(url,"application/json", r )
		if err != nil {
			fmt.Println("Get2: error creating request: ",err.Error())
		}
		fmt.Println("statusCode: ",res.StatusCode )
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}
}

func ReturnData(w http.ResponseWriter, r *http.Request){

	// Listen to tracker

	var getmsg getMsg
	if r.Method == http.MethodPost {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ", err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ", err)
		}
		if err := json.Unmarshal(body, &getmsg); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			fmt.Println(err.Error())
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ", err)
			}
		}

		fmt.Println(getmsg.Key,": "+"node: "+getmsg.NodeID+", "+ getmsg.Name)
		err = ioutil.WriteFile(path + "/src/local/"+getmsg.Key+"/"+getmsg.Name, []byte(getmsg.Text), 0777)
		if err != nil {
			fmt.Println("Peer: error creating/writing file p2p", err.Error())
		}
		numGets++
		if numGets==totalPartsNum-1{
			fmt.Println("ELAPSED: ",time.Since(startGet))
		}


	}
}






















