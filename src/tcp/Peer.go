package tcp


import ("net"
 	"fmt"
 	"bufio"
	"math"
	"strconv"
	"io/ioutil"
	"os"
	"time"
	"sync"

)
//var currentPart int
var peerNum string
var mutex = &sync.Mutex{}
var size int
var totalPartsNum int
var receivedPart int = 0
func PeerListen(port string, peersPorts []string) {
	var err error
	var currentPart int
	// start tracking time
	start:= time.Now()

	// last char of "port" is peerID
	peerNum=strconv.Itoa(int(port[len(port)-1]-'0'))

	finish := make(chan int)

	// listen to other peers
	go p2pListen(":809"+peerNum, finish)

	// obtain peer's port
	selfPort:=port[len(port)-5:]

	//peers := setP2Pconnections(peersPorts, selfPort)

	// Get Peer number from port
	peerNum:=strconv.Itoa(int(port[len(port)-1]-'0'))
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum,07777)
	fmt.Println("Peer listening...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}
	defer ln.Close()
	// It will first read the size of the data to be received,
	// 	then it will change the limit to EOF,
	//	when EOF is reached, the limit will change in order to read next size
	firstMssg:=true

	var partSize int
	var partBuffer []byte

	// accept connection on port
	conn, err := ln.Accept()
	defer conn.Close()
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
			fmt.Println("		about to read")
			// reading partial buffer from connection
			_,err=conn.Read(partBuffer)
			//fmt.Println(conn.)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("		after reading")

			//----------------------------------------------------------------------OPTIONAL--------
			// create new file
			newFileName:= "newFile"+"_"+strconv.Itoa(currentPart)+"_"
			currentPart++ //updating part number, to be used to create new file
			// write / save buffer to file
			err=ioutil.WriteFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum+"/"+newFileName, partBuffer, 0777)
			if err != nil {
				fmt.Println("Peer: error creating/writing file", err.Error())
			}
			fmt.Println("peer: "+peerNum+", currentPart:			 "+strconv.Itoa( currentPart)+", "+string(partBuffer[:15]))
			peers := setP2Pconnections(peersPorts, selfPort)
			sendToPeers(partBuffer, peers, selfPort)
			if (currentPart)==(totalPartsNum/3) {
				fmt.Println("********************************************** WAITING ",peerNum)
				elapsed:= time.Since(start)
				fmt.Println("Peer: "+peerNum+" |ELAPSED= ",elapsed)
				<-finish
				fmt.Println("FINISHED//////////////////////////////////////////")
				return
			}
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

func sendToPeers(partBuffer []byte, peers []peer, selfPort string){
	// Only send what TRACKER sent (do not send what peers sent)
	fmt.Println("sendToPeers start ... ...",selfPort[len(selfPort)-1]-'0')

	for _ , peer := range peers {
		_, err := fmt.Fprintf(peer.conn, string(partBuffer))
		fmt.Println(selfPort+" sending "+string(partBuffer[:15])+" to "+peer.conn.RemoteAddr().String())
		if err != nil {		// receptor finished
			fmt.Println(err.Error())
			fmt.Println("sendToPeers ERROR")
			return
		}
	}
	fmt.Println("sendToPeers end ... ... ",selfPort[len(selfPort)-1]-'0')
}

func setP2Pconnections (peersPorts []string, selfPort string)[]peer{
	var auxPeer peer
	var err	error
	// create list of peers (excluding itself)
	peers := make([]peer,0)
	time.Sleep(1*time.Second)
	for _, peerPort := range peersPorts{
		if peerPort[len(peerPort)-1]  !=  selfPort[len(selfPort)-1]{
			//if last char of ports is the same, { don't add to "peers" }
			auxPeer.port=peerPort
			auxPeer.conn, err = net.Dial("tcp", peerPort)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println(err.Error())
			}
			peers = append(peers, auxPeer)


		}

	}
	return peers
}


func p2pListen(port string, finish chan int){

	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf("___"+err.Error())
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		go handleP2Pconnection(conn, finish)
	}

}

func handleP2Pconnection(conn net.Conn, finish chan int) {
	//var currentPart int
	var partSize int
	var partBuffer []byte
	fmt.Println("			handleP2Pconnection START ...", conn.LocalAddr())

	// sizing buffer to read from connection
	//partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
	partSize = fileChunk
	partBuffer = make([]byte, partSize)

	// reading partial buffer from connection
	n, err := conn.Read(partBuffer)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Println("	//	//	handleP2Pconnection BYTES: ", n)
	//----------------------------------------------------------------------OPTIONAL--------
	// create new file
	mutex.Lock()
	receivedPart++                //updating part number, to be used to create new file
	mutex.Unlock()
	newFileName := "newFile" + "_" + strconv.Itoa(receivedPart) + "_____"


	peerID:=int( conn.LocalAddr().String()[len(conn.LocalAddr().String())-1]-'0')		// peerID = last character of local address
	fmt.Println("ID: ",peerID)
	// write / save buffer to file
	err = ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/chunksToSend" + strconv.Itoa( peerID)+ "/" + newFileName, partBuffer, 0777)
	if err != nil {
		fmt.Println("Peer: error creating/writing file", err.Error())
	}
	//fmt.Println("currentPart (p2p): 			 ", receivedPart)
	//fmt.Println("RECEIVED= ",receivedPart)
	if receivedPart == (totalPartsNum / 2) {
		fmt.Println("	SENDING FINISH ...	--- 		...")
		finish <- 1
	  // send finish message
	}else{fmt.Println("normally exiting HANDLE")}
	//--------------------------------------------------------------------------------------

}



