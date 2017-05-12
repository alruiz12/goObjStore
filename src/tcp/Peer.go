package tcp


import ("net"
 	"fmt"
 	"bufio"
	"math"
	"strconv"
	"io/ioutil"
	"os"
	"time"
	"strings"

)

func PeerListen(port string, peersPorts []string) {
	// listen to other peers
	go p2pListen()

	// connect to other peers
	peers := make([]peer,len(peersPorts-1))
	for index, peerPort := range peersPorts{
		if strings.Compare(peerPort,port[len(port)-5:])!=0{
			//if peerPort == ":self port" { don't add to "peers" }
			peers[index].port=peerPort
		}

		// connect to this socket
		peers[index].conn, err := net.Dial("tcp", peerPort)    //ex:"127.0.0.1:8081"
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Get Peer number from port
	start:= time.Now()
	peerNum:=strconv.Itoa(int(port[len(port)-1]-'0'))
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum,07777)
	fmt.Println("Peer listening...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}
	//defer ln.Close()
	// It will first read the size of the data to be received,
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
	//defer conn.Close()
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
			fmt.Println("Peer: "+peerNum+" = size: ", size)

			if err != nil {
				fmt.Println(err.Error())
			}

			totalPartsNum= int(math.Ceil(float64(size)/float64(fileChunk)))
			currentPart=0
			firstMssg=false
			fmt.Println("Peer: "+peerNum+"|  total parts: ",totalPartsNum)


		}else{ // not first message, read content

			// sizing buffer to read from connection
			partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
			partBuffer=make([]byte,partSize)

			// reading partial buffer from connection
			_,err=conn.Read(partBuffer)
			if err != nil {
				fmt.Println(err)
			}

			//----------------------------------------------------------------------OPTIONAL--------
			// create new file
			newFileName:= "newFile"+"_"+strconv.Itoa(currentPart)+"_"

			currentPart++ //updating part number, to be used to create new file

			// write / save buffer to file
			err=ioutil.WriteFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum+"/"+newFileName, partBuffer, 0777)
			if err != nil {
				fmt.Println("Peer: error creating/writing file", err.Error())
			}
			fmt.Println("currentPart:			 ", currentPart)
			sendToPeers(partBuffer, peersPorts)
			if currentPart==totalPartsNum {
				fmt.Println("Exiting")
				elapsed:= time.Since(start)
				fmt.Println("Peer: "+peerNum+" |ELAPSED= ",elapsed)
				return}
			//--------------------------------------------------------------------------------------

		}

		/*

		// @Todo: check if different from previous = compute hash and compare to keys in map
		newHash:=GetMD5Hash(message)
		vars.FilesMap.Mutex.Lock()
		defer vars.FilesMap.Mutex.Unlock()
		files, exists:=vars.FilesMap.Content[newHash]
		if !exists {

		}




		*/
	}
	time.Sleep(15 * time.Minute)
}

func sendToPeers(partBuffer []byte, peers []peer){
	// @Todo: decide how will peers know each other

	// open connections (and save them in data structure)


	// @Todo: send data
	for peer , index := range peers {
		_, err := fmt.Fprintf(peer.conn, string(partBuffer))
		fmt.Println("to peer: ",index)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func p2pListen(){

}