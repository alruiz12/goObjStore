package tcp
import (
	"net"
	"fmt"
	"bufio"
	//"strings" // only needed below for sample processing
	"os"
	//"github.com/alruiz12/simpleBT/src/vars"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"math"
	"io/ioutil"
	"time"
)
const fileChunk = 1*(1<<10) // 1 KB
func TrackerListen(port string) {

	fmt.Println("Launching tracker...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}

	// It will first read the size of the data to be sent,
	// 	then it will change the limit to EOF,
	//	when EOF is reached, the limit will change in order to read next size
	firstMssg:=true
	var size int
	var totalPartsNum int
	var currentPart int
	var partSize int
	var partBuffer []byte

	// accept connection on port
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		if firstMssg==true {
		// listen for message containing the size of the data to be sent
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
			}
			// remove trailing char ( '/n' )
			message=message[:len(message)-1]

			size, err = strconv.Atoi(message)
			fmt.Println("size: ", size)

			if err != nil {
				fmt.Println(err.Error())
			}
			conn.Write([]byte("OK" + "\n"))
			return
			totalPartsNum= int(math.Ceil(float64(size)/float64(fileChunk)))
			currentPart=0
			partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
			partBuffer=make([]byte,partSize)
			firstMssg=false

		}else{ // not first message, read content
			_, err := bufio.NewReader(conn).Read(partBuffer)
			if err != nil {
				fmt.Println(err.Error())
			}


			//----------------------------------------------------------------------OPTIONAL--------
			// create new file
			newFileName:= "newFile"+"_"+strconv.Itoa(currentPart)+"_"
			_, err =  os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/DoNotDelete/"+newFileName)
			if err != nil {
				fmt.Println(err)
				return
			}
			// write / save buffer to file
			ioutil.WriteFile(newFileName, partBuffer, os.ModeAppend)
			//--------------------------------------------------------------------------------------

		}

		/*

		// output message received
		fmt.Print("Message Received: ", string(message))

		// @Todo: check if different from previous = compute hash and compare to keys in map
		newHash:=GetMD5Hash(message)
		vars.FilesMap.Mutex.Lock()
		defer vars.FilesMap.Mutex.Unlock()
		files, exists:=vars.FilesMap.Content[newHash]
		if !exists {

		}


			//@Todo: open peer connections send metadata to peers,
			//@Todo: chunking algorithm, send chunks
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))

		*/
	}
	fmt.Println(totalPartsNum)
	time.Sleep(1 * time.Minute)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}