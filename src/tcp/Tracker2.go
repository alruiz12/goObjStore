package tcp
import (
	"net"
	"fmt"
	"os"
	"strconv"
	"time"
	"math"
	"bufio"
	"strings"
	"crypto/md5"
	"encoding/hex"

)

const fileChunk = 1*(1<<10) // 1 KB
type peer struct {
	addr	*net.TCPAddr
	conn	*net.TCPConn
}
func TrackerDivideLoad2(IP string, ports []int, filePath string) {
	var status=0
	// 0 = send new size
	// 2 = send content

	go func() {
		var err error
		// create list of (port, connection)
		peers := make([]peer,len(ports))
		for index, port := range ports{
			peers[index].addr, err =net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port) )
			if err != nil {
				fmt.Println("Tracker: ResolveTCPADDr ERROR: ", err.Error())
			}

			// connect to this socket
			peers[index].conn,err =net.DialTCP("tcp4", nil ,peers[index].addr)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println(" DialTCP error	"+err.Error())
			}
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		text:=strconv.FormatInt(fileInfo.Size(),10)	// size
		fmt.Println("SIZE: ",text)
		size,_:=strconv.Atoi(text)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		totalPartsNum:= int(math.Ceil(float64(size)/float64(fileChunk)))
		currentPart := 0
		// ----------------NEW------------------------------
		var currentPeer peer
		var partSize int
		var currentNum int = 0
		var partBuffer []byte
		var n int
		// -------------------------------------------------

		for currentPart<totalPartsNum{
			if status == 2 {	// send content
				currentPeer=peers[currentNum]
				partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
				partBuffer=make([]byte,partSize)
				_,err = file.Read(partBuffer)
				// send to sockets
				n, err = currentPeer.conn.Write(partBuffer)
				if err != nil {
					fmt.Println("ERROR sendig content !!!!!!!! ",err.Error())
				}
				response, err := bufio.NewReader(currentPeer.conn).ReadString('|')
				if err != nil {
					fmt.Println("Error in response from peer ", err.Error())
				}
				//fmt.Println("		RESPONSE: "+response)


				for strings.Compare(string(response[:len(response) - 1]), string(partBuffer[:5])) != 0 {
					n, err = currentPeer.conn.Write(partBuffer)
					if err != nil {
						fmt.Println("ERROR sendig content !!!!!!!! ", err.Error())
					}
					fmt.Println("re send ", string(partBuffer[:5]))
					//time.Sleep(1*time.Second)
					response, err = bufio.NewReader(currentPeer.conn).ReadString('|')
					if err != nil {
						fmt.Println("Error in response from peer ", err.Error())
					}
					fmt.Println("#")
				}



				fmt.Println("Tracker: n= "+ strconv.Itoa( n ) +" currentPart: "+ strconv.Itoa(currentPart) +
					" current peer"+ strconv.Itoa( currentPeer.addr.Port) +"|"+ string(response) +" bytes:"+
					strconv.Itoa(n))

				currentPart++
				currentNum=(currentNum+1)%3
			}else{	// status == 0 -> send size (already stored in 'text')
				// no need to update 'text'
				for _, peer := range peers {
					_, err = fmt.Fprintf(peer.conn, text + string('|'))
					if err != nil {
						fmt.Println("Error sending size")
						fmt.Println(err.Error())
					}
				}
				status=2
			}

		} // <-- for
	}() // <-- go func
	fmt.Println("tracker finishing ....................................")
	time.Sleep(15 * time.Minute)
}
/*
computes the md5 hash for the string given
@param path to the file we want to split
returns the computed hash
*/
func getMd5sum(data []byte) string{
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))

}