package tcp
import ("net"
	"fmt"
	"bufio"
	"strconv"
	"time"
	"io/ioutil"
	"os"
	"math"
	"sync"

)
var mutex = &sync.Mutex{}
var send = &sync.Mutex{}
var read =&sync.Mutex{}
var start= time.Now()
var size int
//var totalPartsNum int

func PeerListen2(IP string, port int, peersPorts []int) {
	var totalPartsNum int

	var err error
	var partBuffer []byte
	finish := make(chan int)
	hashMap := make (map[string]bool)

	// last decimal of "port" is peerID
	peerNum:=strconv.Itoa(port%10)



	// listen to other peers
	//auxAddr, err := net.ResolveTCPAddr("tcp4",":"+"809"+peerNum )
	if err != nil {
		fmt.Println("Error creating auxAddress for p2pListen ",err.Error())
	}




	// resolve receivers (peers) addresses
	addresses := make([]*net.TCPAddr,len(peersPorts))
	for index, nport := range peersPorts{
		addresses[index], err =net.ResolveTCPAddr("tcp4",":"+strconv.Itoa(nport) )
	}



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

	firstMssg:=true


	//var partSize int
	var currentPart int = 0
	// Create folder to save received data
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum,07777)
	var originalMessage string
	for{
		if firstMssg==true {
			originalMessage, err = bufio.NewReader(conn).ReadString('|')
			if err != nil {
				fmt.Println("readString:   ",err.Error())
			}
			// remove trailing char ( '/n' )
			message:=originalMessage[:len(originalMessage)-1]

			size, err = strconv.Atoi(message)
			if size != 0 {
				fmt.Println("------------____------------Peer: " + peerNum + " = size: ", size)
				if err != nil {
					fmt.Println(err.Error())
				}
				firstMssg = false
				partBuffer=[]byte(originalMessage)
				// GO sendToPeers2(partBuffer, peers, port)
				totalPartsNum = int(math.Ceil(float64(size) / float64(fileChunk)))
				//go p2pListen2(auxAddr, finish, totalPartsNum)
				fmt.Println("total parts num ",totalPartsNum)
			}else{fmt.Println("size == 0, re reading")}
			conn.Write([]byte(string(partBuffer[:5])+string('|')))
		}else{
			fmt.Println("Peer: "+peerNum+" about to read part: ",strconv.Itoa( currentPart) )
			// not first message, read content

			// sizing buffer to read from connection
			partBuffer=make([]byte,fileChunk)
			// reading partial buffer from connection
			n ,err:=conn.Read(partBuffer)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("---- Peer: "+peerNum+" AFTER to read part: ",strconv.Itoa( currentPart) )
			newHash:=getMd5sum(partBuffer)
			_, exists := hashMap[newHash]
			if !exists {
				fmt.Println("nooooooooooooooooooooooooo, ", string(partBuffer[:10]) + " peer: ", peerNum)
				hashMap[newHash]=true
				//----------------------------------------------------------------------OPTIONAL--------
				// create new file
				newFileName:= "newFile"+"_"+strconv.Itoa(currentPart)+"_"
				currentPart++ //updating part number, to be used to create new file
				// write / save buffer to file
				err=ioutil.WriteFile(os.Getenv("GOPATH")+"/src/github.com/alruiz12/simpleBT/src/chunksToSend"+peerNum+"/"+newFileName, partBuffer, 0777)
				if err != nil {
					fmt.Println("Peer: error creating/writing file", err.Error())
				}else{fmt.Println("*	"+newFileName+ " "+ peerNum+ "  "+string(partBuffer[:10])+" n= "+strconv.Itoa(n) )}
				//peers := setP2Pconnections(addresses, port)
				//peers:= setP2Pconnections2(addresses, port)
				//time.Sleep(1 * time.Second)
				//go sendToPeers2(partBuffer, peers, port)
				if (currentPart)>(totalPartsNum/3) {
					//conn.Write([]byte(string(partBuffer[:5])+string('|')))
					elapsed:= time.Since(start)
					fmt.Println("Peer: "+peerNum+" |ELAPSED= ",elapsed, " currentPart = ",currentPart, "total parts: ",totalPartsNum)
					fmt.Println("FINISHED")
					<-finish
					fmt.Println("finsh REEEEEEEEEEEECEIVED")
					break

				}


			}
			send.Lock()
			conn.Write([]byte(string(partBuffer[:5])+string('|')))
			send.Unlock()
			fmt.Println("Peer: "+ peerNum +"origin response: ",string(partBuffer[:5])+string('|'))

		}
		//conn.Write([]byte(string(partBuffer[:5])+string('|')))
	}
	elapsed:= time.Since(start)
	fmt.Println("Peer: "+peerNum+" |ELAPSED= ",elapsed, " currentPart = ",currentPart, "total parts: ",totalPartsNum)
	fmt.Println("FINISHED")
	time.Sleep(15 * time.Minute)
}




func p2pListen2(addr *net.TCPAddr, finish chan int, totalPartsNum int){
	var receivedPartOriginal  int =0
	receivedPart := &receivedPartOriginal
	ln, err := net.ListenTCP("tcp4", addr) //ex: ":8081"
	if err!=nil{
		fmt.Println("p2pListen ListenTCP ",err.Error())
	}
	defer ln.Close()
	i:=0
	for (i*3) < (totalPartsNum*2){
		select {
		case <-finish:
		default: conn, err := ln.AcceptTCP()
			if err != nil {
				fmt.Println("conn p2p ",err.Error())
			}
			go handleP2Pconnection2(conn, finish, receivedPart, totalPartsNum )
			i++
		}

	}
	totalTime:=time.Since(start)
	totalTimeStri:=totalTime.String()
	totalTimeStri=totalTimeStri[:len(totalTimeStri)-2]
	time, _:= strconv.ParseFloat(totalTimeStri,64)
	fmt.Println(  time )
	fmt.Println("/////////////////////////////////////////////////////////////////////////////////")
}


func handleP2Pconnection2(conn *net.TCPConn, finish chan int, receivedPart *int, totalPartsNum int) {
	buf := make([]byte, fileChunk)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("handleP2Pconnection2: ERROR ", err.Error() )
	}
	defer conn.Close()

	mutex.Lock()
	*receivedPart++                //updating part number, to be used to create new file
	mutex.Unlock()
	fmt.Println("receivedPart-----------------------: ",*receivedPart)


	//----------------------------------------------------------------------OPTIONAL--------
	// create new file
	newFileName := "newFile" + "_" + strconv.Itoa(*receivedPart) + "_____"
	peerID:=int( conn.LocalAddr().String()[len(conn.LocalAddr().String())-1]-'0')		// peerID = last character of local address
	fmt.Println("handlep2pconn ID: ",peerID)
	// write / save buffer to file
	err = ioutil.WriteFile(os.Getenv("GOPATH") + "/src/github.com/alruiz12/simpleBT/src/chunksToSend" + strconv.Itoa( peerID)+ "/" + newFileName, buf, 0777)
	if err != nil {
		fmt.Println("Peer: error creating/writing file", err.Error())
	}
	//--------------------------------------------------------------------------------------

}



func setP2Pconnections2 (addresses []*net.TCPAddr, selfPort int)[]peer{
	var auxPeer peer
	var err	error
	// create list of peers (excluding itself)
	peers := make([]peer,0)
	for _, address := range addresses {
		if address.Port%10  !=  selfPort%10{
			//if last char of ports is the same, { don't add to "peers" }
			auxPeer.addr=address
			auxPeer.conn, err = net.DialTCP("tcp",nil, address)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println("setP2P error:",err.Error())
			}
			peers = append(peers, auxPeer)


		}
	}
	return peers
}

func sendToPeers2(buf []byte, peers []peer, selfPort int){
	// Only send what TRACKER sent (do not send what peers sent)
	fmt.Println("SEEEEEEEEEEEEEEEEEEEEEEEND")
	var err error
	for _ , peer := range peers {
		mutex.Lock()
		_, err = peer.conn.Write(buf)
		mutex.Unlock()

		if err != nil {		// receptor finished
			fmt.Println(err.Error())
			fmt.Println("sendToPeers ERROR ",err.Error())
			return
		}
	}
}

