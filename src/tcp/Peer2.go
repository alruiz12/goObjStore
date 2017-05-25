package tcp
import ("net"
	"fmt"
	"bufio"
	"strconv"
	"time"
	"io/ioutil"
	"os"
	"math"

)

func PeerListen2(IP string, port int, peersPorts []int) {
	start:= time.Now()
	var err error
	var size int
	var partBuffer []byte
	finish := make(chan int)
	hashMap := make (map[string]bool)
	// last decimal of "port" is peerID
	peerNum:=strconv.Itoa(port%10)


/*
	// listen to other peers
	auxAddr, err := net.ResolveTCPAddr("tcp4",":"+"809"+peerNum )
	if err != nil {
		fmt.Println("Error creating auxAddress for p2pListen ",err.Error())
	}
	go p2pListen2(auxAddr, finish)
*/
/*

	// resolve receivers (peers) addresses
	addresses := make([]*net.TCPAddr,len(peersPorts))
	for index, nport := range peersPorts{
		addresses[index], err =net.ResolveTCPAddr("tcp4",":"+strconv.Itoa(nport) )
	}
	peers:= setP2Pconnections2(addresses, port)
*/

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
	for {
		if firstMssg==true {
			originalMessage, err = bufio.NewReader(conn).ReadString('|')
			if err != nil {
				fmt.Println("readString:   ",err.Error())
			}
			// remove trailing char ( '/n' )
			//fmt.Println("originalMessage: ", originalMessage)
			message:=originalMessage[:len(originalMessage)-1]

			size, err = strconv.Atoi(message)
			if size != 0 {
				fmt.Println("------------____------------Peer: " + peerNum + " = size: ", size)
				if err != nil {
					fmt.Println(err.Error())
				}
				firstMssg = false
				//fmt.Println("Peer: " + peerNum + " received size!!")
				partBuffer=[]byte(originalMessage)
				//fmt.Println("original msg after sending: ",originalMessage," partBuffer: "+string(partBuffer) )
				// GO sendToPeers2(partBuffer, peers, port)
				totalPartsNum = int(math.Ceil(float64(size) / float64(fileChunk)))
				fmt.Println("total parts num ",totalPartsNum)
				//time.Sleep(100 * time.Millisecond)
			}else{fmt.Println("size == 0, re reading")}

		}else{
			fmt.Println("Peer: "+peerNum+" about to read part: ",strconv.Itoa( currentPart) )
			// not first message, read content

			// sizing buffer to read from connection
			//partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
			partBuffer=make([]byte,1024)
			// reading partial buffer from connection
			n ,err:=conn.Read(partBuffer)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("---- Peer: "+peerNum+" AFTER to read part: ",strconv.Itoa( currentPart) )
			newHash:=getMd5sum(partBuffer)
			_, exists := hashMap[newHash]
			if !exists {
				fmt.Println("////////////////// DOES NOT EXIST")
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
				//sendToPeers(partBuffer, peers, port)

				if (currentPart)==(totalPartsNum/3) {
					conn.Write([]byte(string(partBuffer[:5])+string('|')))
					elapsed:= time.Since(start)
					fmt.Println("Peer: "+peerNum+" |ELAPSED= ",elapsed, " currentPart = ",currentPart, "total parts: ",totalPartsNum)
					<-finish
					fmt.Println("FINISHED")
					return
				}else{
					fmt.Println("CURRENTPART=",currentPart," TOTAL PARTS=",totalPartsNum)
				}

			} else {fmt.Println("/////EXISTS ALREADY!!")}

		}
		conn.Write([]byte(string(partBuffer[:5])+string('|')))
	}
	time.Sleep(15 * time.Minute)
}




func p2pListen2(addr *net.TCPAddr, finish chan int){
	//fmt.Println("	p2plisten start")
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

		go handleP2Pconnection2(conn, finish)
	}
}


func handleP2Pconnection2(conn *net.TCPConn, finish chan int) {
	//fmt.Println("		handle p2p2 connection start")
	//buf := make([]byte, partSize)
	_, err := bufio.NewReader(conn).ReadString('|')
	if err != nil {
		fmt.Println("readString:   ",err.Error())
	}
	/*_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("handleP2Pconnection2: ERROR ", err.Error() )
	}*/
	//fmt.Println("		content: "+originalMessage)
}



func setP2Pconnections2 (addresses []*net.TCPAddr, selfPort int)[]peer{
	var auxPeer peer
	var err	error
	// create list of peers (excluding itself)
	peers := make([]peer,0)
	time.Sleep(1*time.Second)
	for _, address := range addresses {
		if address.Port%10  !=  selfPort%10{
			//if last char of ports is the same, { don't add to "peers" }
			auxPeer.addr=address
			auxPeer.conn, err = net.DialTCP("tcp",nil, address)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println(err.Error())
			}
			peers = append(peers, auxPeer)


		}
	}
	return peers
}

func sendToPeers2(buf []byte, peers []peer, selfPort int){
	// Only send what TRACKER sent (do not send what peers sent)
	//fmt.Println("sendToPeers start ... ...", selfPort%10)
	var err error
	for _ , peer := range peers {
		_, err = peer.conn.Write(buf)
		//fmt.Println(strconv.Itoa( selfPort) +" sending "+string(buf)+" with a length: "+ strconv.Itoa(len(buf)) +" to "+peer.conn.RemoteAddr().String())
		if err != nil {		// receptor finished
			fmt.Println(err.Error())
			fmt.Println("sendToPeers ERROR ",err.Error())
			return
		}
	}
	//fmt.Println("sendToPeers end ... ... ",selfPort%10)
}

