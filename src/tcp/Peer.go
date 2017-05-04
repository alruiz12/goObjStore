package tcp


import "net"
import "fmt"
import "bufio"
import (
	"math"
	"strconv"
	"io/ioutil"
	"os"
	"time"
)


/*
func getNextFilePart(file string) string{

}
*/
func PeerListen(port string) {
	peerNum:=strconv.Itoa(int(port[len(port)-1]-'0'))
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum,07777)
	fmt.Println("Peer listening...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}

	// It will first read the size of the data to be received,
	// 	then it will change the limit to EOF,
	//	when EOF is reached, the limit will change in order to read next size
	firstMssg:=true
	var size int
	var totalPartsNum int
	var currentPart int
	var partSize int
	var partBuffer []byte
	//var err error

	// accept connection on port
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		if firstMssg==true {
			// listen for message containing the size of the data to be received
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
			totalPartsNum= int(math.Ceil(float64(size)/float64(fileChunk)))
			currentPart=0
			firstMssg=false
			fmt.Println("total parts: ",totalPartsNum)


		}else{ // not first message, read content
			partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
			partBuffer=make([]byte,partSize)
			/*
			if currentPart==0 {
				reader = bufio.NewReader(conn)
			} */
			fmt.Println("		not first message")
			_,err=conn.Read(partBuffer)
			if err != nil {
				fmt.Println(err.Error())
			}



			//----------------------------------------------------------------------OPTIONAL--------
			// create new file
			newFileName:= "newFile"+"_"+strconv.Itoa(currentPart)+"_"
			/*
			_, err =  os.Create(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend/"+newFileName)
			if err != nil {
				fmt.Println(err)
				return
			}

			file, err:=os.Open(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend/"+newFileName)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			*/

			currentPart++
			// write / save buffer to file
			//file.Write(partBuffer[:n])
			ioutil.WriteFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum+"/"+newFileName, partBuffer, 0777)
			if currentPart==totalPartsNum {
				fmt.Println("Exiting")
				return}
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
	time.Sleep(5 * time.Minute)
}