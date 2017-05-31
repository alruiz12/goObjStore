package httpGo
import(
	"os"
	"fmt"
	"strconv"
	"math"
	"net/http"
	"io"
	"bytes"
	"mime/multipart"
	"time"
	"encoding/json"
)

const fileChunk = 1*(1<<10) // 1 KB
//const fileChunk = 8*(1<<20) // 8 MB
func TrackerDivideLoad(filePath string, addr string, peers []string){
	time.Sleep(5 * time.Second)
	var currentPart int = 0
	var currentPeer string
	var partSize int
	var currentNum int = 0
	var partBuffer []byte
	var err error
	var peerURL string
	var body *bytes.Buffer
	var writer *multipart.Writer
	var fileName string

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

	totalPartsNum:= int(math.Ceil(float64(size)/float64(fileChunk)))

	for currentPart<totalPartsNum{
		currentPeer=peers[currentNum]
		partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
		partBuffer=make([]byte,partSize)
		_,err = file.Read(partBuffer)	// Get chunk

		peerURL = "http://"+currentPeer+"/PeerListen"

		body = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
		fileName = strconv.Itoa(currentPart)
		part, err := writer.CreateFormFile("file", fileName) // Todo: declare
		if err != nil {
			fmt.Println("creating Form file")
			fmt.Println(err)
		}
		_, err = io.Copy(part, file)
		err=writer.Close()
		if err != nil {
			fmt.Println(err)
		}
		request, err := http.NewRequest("POST", peerURL, body) // Todo: declare
		if err != nil {
			fmt.Println("Error creating request = ",err)
		}
		request.Header.Set("Content-Type", writer.FormDataContentType())
		//_, err = http.DefaultClient.Do(request) // Todo: declare
		if err != nil {
			fmt.Println("Error doing requewst = ",err)
		}
		/*
		if res.StatusCode != 200 {
			fmt.Println("Success expected: %d", res.StatusCode)
		}

*/		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode("AAAAAA")
		http.Post(peerURL, "application/json", &buf )

		return
		//currentNum++
	}




	// Get chunk

	// Send to current peer

}