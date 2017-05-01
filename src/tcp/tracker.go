package tcp


import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func TrackerListen(port string) {

	fmt.Println("Launching tracker...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", port) //ex: ":8081"
	if err!=nil{
		fmt.Printf(err.Error())
	}

	// accept connection on port
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err.Error())
	}
	go func() {
		for {
			// @Todo limit buffer ('\n' not valid)
			// will listen for message to process ending in newline (\n)
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
			}
			// output message received
			fmt.Print("Message Received:", string(message))

			// @Todo address peers: check if different from previous,
				//@Todo: open peer connections send metadata to peers,
				//@Todo: chunking algorithm, send chunks
			// sample process for string received
			newmessage := strings.ToUpper(message)
			// send new string back to client
			conn.Write([]byte(newmessage + "\n"))
		}
	}()
}