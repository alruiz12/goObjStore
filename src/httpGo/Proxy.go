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
	"sync"
	"math/rand"
	"github.com/alruiz12/simpleBT/src/conf"
)

const fileChunk = 1*(1<<10) // 1 KB
//const fileChunk = 8*(1<<20) // 8 MB

type msg struct {
	NodeList [][]string
	Num int
	Hash string
	Text []byte 
	CurrentNode int
	Name int
}
type hashMsg struct {
	Hash string
}
var totalPartsNum int
var startGet time.Time
var numGets int = 0
func PutObjProxy(filePath string, trackerAddr string, numNodes int, putOK chan bool) {

	time.Sleep(1 * time.Second)
	var hash string = md5sum(filePath)
	var err error

	// Aask tracker for nodes
	requestJson := `{"Quantity":"` + strconv.Itoa(numNodes) + `","ID":"` + hash + `","Type":"object"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + trackerAddr + "/GetNodes"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Put: error creating request: ", err.Error())
		putOK <- false
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Put: error sending request: ", err.Error())
		putOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		putOK <- false
		return
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		putOK <- false
		return
	}
	var nodeList [][]string
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Put: error unprocessable entity: ", err.Error())
		putOK <- false
		return
	}
	if err != nil {
		fmt.Println("Put: error receiving response: ", err.Error())
		putOK <- false
		return
	}
	var currentPart int = 0
	var partSize int
	var currentNum int = 0
	var partBuffer []byte
	var writer *multipart.Writer
	var buf bytes.Buffer
	_, _ = writer, buf // avoiding declared but not used

	var auxList []bool
	var i int = 0
	for i < len(nodeList) {
		auxList = append(auxList, false)
		i++
	}
	httpVar.DirMutex.Lock()
	httpVar.HashMap[hash] = auxList
	httpVar.DirMutex.Unlock()



	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		putOK <- false
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	text := strconv.FormatInt(fileInfo.Size(), 10)        // size
	size, _ := strconv.Atoi(text)
	if err != nil {
		fmt.Println(err.Error())
		putOK <- false
		return
	}
	httpVar.TotalNumMutex.Lock()
	totalPartsNum = int(math.Ceil(float64(size) / float64(fileChunk)))
	httpVar.TotalNumMutex.Unlock()
	var currentAdr int = 0

	for currentNum < len(nodeList) {
		rpipe, wpipe := io.Pipe()
		mHash := hashMsg{Hash:hash}
		go func() {
			// save buffer to object
			err = json.NewEncoder(wpipe).Encode(mHash)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
				putOK <- false
				return		// Todo: does this exit only goRoutine?
			}
			defer wpipe.Close()                     // close pipe //when go routine finishes
		}()

		// Prepare node for content
		_, err = http.Post("http://" + nodeList[currentNum][currentAdr] + "/prepSN", "application/json", rpipe)
		if err != nil {
			fmt.Println("to prepSN, Error sending http POST ", err.Error())
			putOK <- false
			return
		}
		currentNum++
	}
	currentNum = 0
	//return
	var wg sync.WaitGroup
	wg.Add(totalPartsNum)
	for currentPart < totalPartsNum {
		partSize = int(math.Min(fileChunk, float64(size - (currentPart * fileChunk))))
		partBuffer = make([]byte, partSize)
		_, err = file.Read(partBuffer)                // Get chunk
		m := msg{NodeList:nodeList, Num:len(nodeList), Hash:hash, Text:partBuffer, CurrentNode:currentNum, Name: currentPart}
		//r, w :=io.Pipe()			// create pipe
		go func(m2 msg, url string) {
			 httpVar.SendReady <- 1
			r, w := io.Pipe()
			go func() {
				// save buffer to object
				err = json.NewEncoder(w).Encode(m2)
				if err != nil {
					fmt.Println("Error encoding to pipe ", err.Error())
					putOK <- false
					return
				}
				defer w.Close()                        // close pipe //when go routine finishes
			}()
			_, err := http.Post(url, "application/json", r)
			if err != nil {
				fmt.Println("Error sending http POST ", err.Error())
				putOK <- false
				return
			}
			defer wg.Done()
			 <-httpVar.SendReady
		}(m, "http://" + nodeList[currentNum][currentAdr] + "/SNPutObj")

		currentPart++
		currentNum = (currentNum + 1) % len(nodeList)

		// Every 'numNodes' iterations, send chunk to next address, first send to different nodes, then change address
		if currentNum == 0 {
			currentAdr = (currentAdr + 1) % len(nodeList[currentNum])
		}
	}
	wg.Wait()
	putOK <- true
	//fmt.Println("Proxy's routines finished")
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


func GetObjProxy(Key string, proxyAddr []string, trackerAddr string){
	time.Sleep(1 * time.Second)
//	var chunk msg
	// Ask tracker for nodes
	startGet=time.Now()
	var err error
	// ask tracker for nodes for a given key
	requestJson := `{"ID":"`+Key+`","Type":"object"}`
	reader := strings.NewReader(requestJson)
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
	var nodeList [][]string
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Get: error unprocessable entity: ",err.Error())
		return
	}
	if err != nil {
		fmt.Println("Get: error reciving response: ",err.Error())
	}
	// Create folder for receiving
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/local",+0777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/local/"+Key,0777)

	//var currentAddr int = rand.Intn(len(chunk.NodeList)) @Todo
	var currentAddr int = 0
	// For each node ask for all their Proxy-pieces
	for index, node := range nodeList {
		r, w :=io.Pipe()			// create pipe
		k:=jsonKeyURL{Key:Key, URL:proxyAddr[index]+"/ReturnObjProxy"}

		go func() {
			defer w.Close()			// close pipe when go routine finishes
			// save buffer to object
			err=json.NewEncoder(w).Encode(&k)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
			}
		}()
		url:="http://"+node[currentAddr]+"/SNPutObjGetChunks"
		res, err := http.Post(url,"application/json", r )
		if err != nil {
			fmt.Println("Get2: error creating request: ",err.Error())
		}
		//fmt.Println("statusCode: ",res.StatusCode )
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}

}

func ReturnObjProxy(w http.ResponseWriter, r *http.Request){

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
		httpVar.GetMutex.Lock()
                numGets++
                httpVar.GetMutex.Unlock()

		err = ioutil.WriteFile(path + "/src/local/"+getmsg.Key+"/"+getmsg.Name, getmsg.Text, 0777)
		if err != nil {
			fmt.Println("Peer: error creating/writing file p2p", err.Error())
		}
		httpVar.TotalNumMutex.Lock()

		if numGets==totalPartsNum{
			fmt.Println("GET: ",time.Since(startGet))
		}
		httpVar.TotalNumMutex.Unlock()


	}
}


/*
CheckPieces walks through the subfiles directory, creates a new file to be filled out with the content of each subfile,
and compares the new hash with the original one.
@param path to the file we want to split
Returns true if both hash are identic and false if not
*/
func CheckPiecesObj(key string ,fileName string, filePath string, numNodes int) bool{
	 file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	text := strconv.FormatInt(fileInfo.Size(), 10)        // size
	size, _ := strconv.Atoi(text)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	totalPartsNumOriginal := int(math.Ceil(float64(size) / float64(fileChunk)))

	// Walking through StorageNodes data
	// Subfiles directory
	path:=os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/data/"+key+"/"
	subDir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	currentDir:=0
	for currentDir<numNodes{

		// Create new file
		_, err = os.Create(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src" + fileName+strconv.Itoa(currentDir))
		newFile, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src" + fileName+strconv.Itoa(currentDir), os.O_APPEND | os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer newFile.Close()

		files, err := ioutil.ReadDir(path+subDir[currentDir].Name() )
		//fmt.Println("Entering ", path+subDir[currentDir].Name())
		// Trying to fill out the new file using subfiles (in order)
		var inOrderCount = 0
		var maxTimes int = 0
		var fileNameOriginal= fileName[:len(fileName)-4]

		for inOrderCount<totalPartsNumOriginal {
			for _, file := range files {
				if strings.Compare(file.Name(), fileNameOriginal + strconv.Itoa(inOrderCount)) == 0 || strings.Compare(file.Name(), "P2P" + strconv.Itoa(inOrderCount)) == 0{
					inOrderCount++
					//				fmt.Println(file.Name())
					currentFile, err := os.Open(path + subDir[currentDir].Name() +"/"+ file.Name())
					if err != nil {
						fmt.Println(err)
						return false
					}

					bytesCurrentFile, err := ioutil.ReadFile(path + subDir[currentDir].Name()+"/" +file.Name())

					_, err = newFile.WriteString(string(bytesCurrentFile))
					if err != nil {
						fmt.Println(err)
						return false
					}

					currentFile.Close()
				}

			}
			if inOrderCount == 0 {
				maxTimes++
			}
			if maxTimes > 1 {
				fmt.Println("maxTimes > 1 when looking for ", fileNameOriginal + strconv.Itoa(inOrderCount))
				return false
			}
		}

		// Compute and compare new hash
		newHash := md5sum(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src" + fileName+strconv.Itoa(currentDir))
		//fmt.Println("/src/github.com/alruiz12/simpleBT/src" + fileName+strconv.Itoa(currentDir) , newHash)
		if strings.Compare(key, newHash) != 0 {
			return false
		}



		currentDir++
	}
	if currentDir==0{return false}	// Never got in loop
	//return true

	// Checking Get output (locally)
	path=os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/local/"+key+"/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Create new file
	_, err = os.Create(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src" + fileName)
	newFile, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src" + fileName, os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer newFile.Close()



	// Trying to fill out the new file using subfiles (in order)
	var inOrderCount = 0
	var maxTimes int = 0
	var fileNameOriginal= fileName[:len(fileName)-4]
	for inOrderCount<totalPartsNumOriginal {
		for _, file := range files {
			if strings.Compare(file.Name(), fileNameOriginal + strconv.Itoa(inOrderCount)) == 0 {
				inOrderCount++

				currentFile, err := os.Open(path + file.Name())
				if err != nil {
					fmt.Println(err)
					return false
				}

				bytesCurrentFile, err := ioutil.ReadFile(path + file.Name())

				_, err = newFile.WriteString(string(bytesCurrentFile))
				if err != nil {
					fmt.Println(err)
					return false
				}

				currentFile.Close()
			}

		}
		if inOrderCount == 0 {
			maxTimes++
		}
		if maxTimes > 1 {
			fmt.Println("maxTimes > 1 when looking for ", fileNameOriginal + strconv.Itoa(inOrderCount))
			return false
		}
	}

	// Compute and compare new hash
	newHash := md5sum(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src" + fileName)
	//fmt.Println("/src/github.com/alruiz12/simpleBT/src" + fileName ,newHash)
	if strings.Compare(key, newHash) != 0 {
		return false
	}

	return true
}

type AccInfo struct {
	NodeList [][]string
	Num int
	CurrentNode int
	Name string
}

func CreateAccountProxy(name string, createOK chan bool){
	var nodeList [][]string
	var err error

	// Aask tracker for nodes
	requestJson := `{"Quantity":"` + strconv.Itoa(conf.NumNodes) + `","ID":"` + name + `","Type":"account"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + conf.TrackerAddr + "/GetNodes"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Put: error creating request: ", err.Error())
		createOK <- false
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Put: error sending request: ", err.Error())
		createOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		createOK <- false
		return
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		createOK <- false
		return
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Put: error unprocessable entity: ", err.Error())
		createOK <- false
		return
	}
	if err != nil {
		fmt.Println("Put: error receiving response: ", err.Error())
		createOK <- false
		return
	}

	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{NodeList:nodeList, Num:len(nodeList), CurrentNode:currentPeer, Name:name }

	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			fmt.Println("Error encoding to pipe ", err.Error())
			createOK <- false
			return
		}
		defer w.Close()                        // close pipe //when go routine finishes
	}()
	fmt.Println(nodeList[currentPeer][currentPeerAddr])
	_, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/SNPutAcc", "application/json", r)
	if err != nil {
		fmt.Println("Error sending http POST ", err.Error())
		createOK <- false
		return
	}

	createOK <- true

}







