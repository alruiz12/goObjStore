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
	"io/ioutil"
	"time"
	"strings"
)
const fileChunk = 1*(1<<10) // 1 KB
type peer struct {
	port	string
	conn	net.Conn
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}


func TrackerFile(IP string, ports []string, filePath string) {
	var status=0	// 0 = send new size
	// 1 = resend message
	// 2 = send content

	var allBytes []byte
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
		peers := make([]peer,len(ports))
		for index, port := range ports{
			peers[index].port=port

			// connect to this socket
			peers[index].conn,err =net.Dial("tcp", IP+port)    //ex:"127.0.0.1:8081"
			if err != nil {
				fmt.Println(err.Error())
			}
		}


		text:=strconv.FormatInt(fileInfo.Size(),10)	// size
		for {
			if status == 2 {	// send content
				allBytes,err =ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				text=string(allBytes)
				fmt.Println("sending file data")
			}

			// send to sockets

			if status==2{
				for index, peer := range peers {
					n, err = fmt.Fprintf(peer.conn, text)
					if err != nil {
						fmt.Println(err.Error())
					}
					fmt.Println("peer"+strconv.Itoa(index)+": bytes send: ", n)
				}
			}else{
				for index, peer := range peers {
					n, err = fmt.Fprintf(peer.conn, text + string('\n'))
					if err != nil {
						fmt.Println(err.Error())
					}
					fmt.Println("peer"+strconv.Itoa(index)+": bytes send: ", n)

				}
				status=1
			}


			if status==2 {return }

			// listen for reply
			for index, peer := range peers {
				go func() {
					message, _ := bufio.NewReader(peer.conn).ReadString('\n')
					fmt.Print("Message from peer" + strconv.Itoa(index) + ": " + message)

					if strings.Compare(message, "OK\n") == 0 {
						fmt.Println("Received OK")
						status = 2
					} else {
						status = 1 // resend message
					}
				}()
			}
		}
	}()
	time.Sleep(5 * time.Minute)
}