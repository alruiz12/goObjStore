package tcp
import (
	"net"
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
	"math"
)
const fileChunk = 1*(1<<10) // 1 KB
//const fileChunk = 8*(1<<20) // 8 MB
type peer struct {
	addr	*net.TCPAddr
	conn	*net.TCPConn
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}


func TrackerDivideLoad(IP1 byte,IP2 byte,IP3 byte,IP4 byte, ports []int, filePath string) {
	var status=0	// 0 = send new size
	// 1 = resend message
	// 2 = send content

	var n int
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	go func() {
		// create list of (port, connection)
		fmt.Println("len ports", len(ports))
		peers := make([]peer,len(ports))
		for index, port := range ports{
			fmt.Println("index=", index)
			peers[index].addr.IP =net.IPv4(IP1, IP2, IP3, IP4)
			peers[index].addr.Port=port

			// connect to this socket
			peers[index].conn,err =net.DialTCP("tcp", nil ,peers[index].addr)    //ex:"127.0.0.1:8081"
			if err != nil {
				//fmt.Println("Error creating connection in 'peers', port= "+IP+peers[index].port+" index= "+strconv.Itoa( index))
				fmt.Println("	"+err.Error())
			}
		}


		text:=strconv.FormatInt(fileInfo.Size(),10)	// size
		size,_:=strconv.Atoi(text)
		//allBytes,err =ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		totalPartsNum:= int(math.Ceil(float64(size)/float64(fileChunk)))
		var partSize int
		var partBuffer []byte
		currentPart := 0
		//var content string
		var currentPeer peer
		var currentNum int=0
		for currentPart<totalPartsNum{
			if status == 2 {	// send content
				currentPeer=peers[currentNum]
				partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
				partBuffer=make([]byte,partSize)
				_,err = file.Read(partBuffer)
				//content=string(partBuffer)
				//fmt.Println("sending file data")

				// send to sockets
				n, err= currentPeer.conn.Write(partBuffer)
				//n, err = fmt.Fprintf(currentPeer.conn, content)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println("Tracker: n= ",n/*, content[:25]*/)

				currentPart++
				currentNum=(currentNum+1)%3

			}else{	// status == 0 -> send size (already stored in 'text')
				// OR status ==1 -> re send same content (already stored in 'text')
				// no need to update 'text'
				for index, peer := range peers {
					fmt.Println("index ", index)
					text=text + string('\n')
					copy(partBuffer,text)
					n, err= peer.conn.Write(partBuffer)
					//_, err = fmt.Fprintf(peer.conn, text + string('\n'))
					if err != nil {
						fmt.Println("Error sending size")
						fmt.Println(err.Error())
						//panic(err)
					}

				}
				//status=1
				status=2
			}


		} // <-- for
	}() // <-- go func
	fmt.Println("tracker finishing ....................................")
	time.Sleep(15 * time.Minute)
}