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
	//"encoding/json"
	//"strings"

)


var mutex = &sync.Mutex{}
var size int
var totalPartsNum int
var receivedPart int = 0
func PeerListen(IP string, port int, peersPorts []int) {
	var err error
	var currentPart int
	var partSize int
	var partBuffer []byte
	start:= time.Now()
	finish := make(chan int)

	// last decimal of "port" is peerID
	peerNum:=strconv.Itoa(port%10)



	// listen to other peers
	auxAddr, err := net.ResolveTCPAddr("tcp4",":"+"809"+peerNum )
	if err != nil {
		fmt.Println("Error creating auxAddress for p2pListen ",err.Error())
	}
	go p2pListen(auxAddr, finish)



	addresses := make([]*net.TCPAddr,len(peersPorts))
	for index, nport := range peersPorts{
		addresses[index], err =net.ResolveTCPAddr("tcp4",":"+strconv.Itoa(nport) )
	}

	// Create folder to save received data
	//os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum,07777)

	// listen on all interfaces
	addr, err :=net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port) )
	ln, err := net.ListenTCP("tcp4", addr)//ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}
	defer ln.Close()

	// accept connection on port
	conn, err := ln.AcceptTCP()
	defer conn.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	/*
	var hashList []string
	message, err := bufio.NewReader(conn).ReadString('\n')
	dec := json.NewDecoder(strings.NewReader(message))
	err = dec.Decode(&hashList)
	if err != nil {
		fmt.Println(" error reading list ",err.Error())
	}
	fmt.Println(hashList)
	*/
	firstMssg:=true
	for {
		if firstMssg==true {
			message, err := bufio.NewReader(conn).ReadString('|')
			if err != nil {
				fmt.Println("readString:   ",err.Error())
			}
			// remove trailing char ( '/n' )
			fmt.Println("messagez: ", message)
			message=message[:len(message)-1]

			size, err = strconv.Atoi(message)
			if size != 0 {

				fmt.Println("------------____------------Peer: " + peerNum + " = size: ", size)
				if err != nil {
					fmt.Println(err.Error())
				}
				/*
				totalPartsNum = int(math.Ceil(float64(size) / float64(fileChunk)))
				currentPart = 0*/
				firstMssg = false
				fmt.Println("Peer: " + peerNum + " received size!!")
				//fmt.Println("Peer: " + peerNum + "|  total parts: ", totalPartsNum)

			}else{fmt.Println("size == 0, re reading")}

		}else{
			fmt.Println("Peer "+peerNum+ " waiting ...")
			time.Sleep(5*time.Minute)
			// not first message, read content
			fmt.Println("else: size ", size)
			// sizing buffer to read from connection
			partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
			partBuffer=make([]byte,partSize)
			// reading partial buffer from connection
			n ,err:=conn.Read(partBuffer)
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
			}else{fmt.Println("	*	*	*	"+newFileName+ " "+ peerNum+ "  "+string(partBuffer)+" n= "+strconv.Itoa(n) )}
			peers := setP2Pconnections(addresses, port)
			sendToPeers(partBuffer, peers, port)

			if (currentPart-1)==(totalPartsNum/3) {
				elapsed:= time.Since(start)
				fmt.Println("Peer: "+peerNum+" |ELAPSED= ",elapsed)
				<-finish
				fmt.Println("FINISHED")
				return
			}
			//--------------------------------------------------------------------------------------

		}

	}
	time.Sleep(15 * time.Minute)
}

func sendToPeers(partBuffer []byte, peers []peer, selfPort int){
	// Only send what TRACKER sent (do not send what peers sent)
	fmt.Println("sendToPeers start ... ...", selfPort%10)
	var err error
	for _ , peer := range peers {
		_, err = peer.conn.Write(partBuffer)
		fmt.Println(strconv.Itoa( selfPort) +" sending "+string(partBuffer[:len(partBuffer)/100])+" to "+peer.conn.RemoteAddr().String())
		if err != nil {		// receptor finished
			fmt.Println(err.Error())
			fmt.Println("sendToPeers ERROR ",err.Error())
			return
		}
	}
	fmt.Println("sendToPeers end ... ... ",selfPort%10)
}

func setP2Pconnections (addresses []*net.TCPAddr, selfPort int)[]peer{
	var auxPeer peer
	var err	error
	// create list of peers (excluding itself)
	peers := make([]peer,0)
	time.Sleep(1*time.Second)
	for _, address := range addresses {
		if address.Port%10  !=  selfPort%10{
			//if last char of ports is the same, { don't add to "peers" }
			auxPeer.addr=address
			fmt.Println("SET", address)
			auxPeer.conn, err = net.DialTCP("tcp",nil, address)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println(err.Error())
			}
			peers = append(peers, auxPeer)


		}
	}
	return peers
}


func p2pListen(addr *net.TCPAddr, finish chan int){
	fmt.Println("p2plisten start")
	ln, err := net.ListenTCP("tcp4", addr) //ex: ":8081"
	if err!=nil{
		fmt.Println("p2pListen ListenTCP ",err.Error())
	}
	defer ln.Close()
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println("conn p2p ",err.Error())
		}
		go handleP2Pconnection(conn, finish)
	}
}

func handleP2Pconnection(conn *net.TCPConn, finish chan int) {
	var partSize int
	var partBuffer []byte
	fmt.Println("handleP2Pconnection START ..." )

	// sizing buffer to read from connection
	partSize = fileChunk
	partBuffer = make([]byte, partSize)

	// reading partial buffer from connection
	n, err := conn.Read(partBuffer)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Println("handleP2Pconnection BYTES: ", n)

	mutex.Lock()
	receivedPart++                //updating part number, to be used to create new file
	mutex.Unlock()
	//----------------------------------------------------------------------OPTIONAL--------
	// create new file
	newFileName := "newFile" + "_" + strconv.Itoa(receivedPart) + "_____"
	peerID:=int( conn.LocalAddr().String()[len(conn.LocalAddr().String())-1]-'0')		// peerID = last character of local address
	fmt.Println("handlep2pconn ID: ",peerID)
	// write / save buffer to file
	err = ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/chunksToSend" + strconv.Itoa( peerID)+ "/" + newFileName, partBuffer, 0777)
	if err != nil {
		fmt.Println("Peer: error creating/writing file", err.Error())
	}
	//--------------------------------------------------------------------------------------
	if receivedPart == (totalPartsNum / 2) {
		fmt.Println("	SENDING FINISH ...	--- 		...")
		finish <- 1
	  // send finish message
	}else{fmt.Println("normally exiting HANDLE")}


}



