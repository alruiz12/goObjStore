package conf

import (
	//"strings"
	"errors"
	"fmt"
	//"os"
	"encoding/json"
	"os"
)
//import "fmt"

var currentId int

var torrents Torrents

func getPeersRepo(t Torrent) Peers{
	return t.Peers
}

func addTorrentRepo(t Torrent) Torrent{
	var err error
	currentId += 1
	t.Id = currentId
	torrents = append(torrents, t)
	f, err:=os.OpenFile("torrentsFile",os.O_APPEND|os.O_WRONLY,0666)
	if err!=nil {panic (err)}

	out,err:=json.Marshal(t)
	writtenBytes, err := f.WriteString(string(out) )
	fmt.Println("wrote %d bytes",writtenBytes)
	f.Close()
	return t
}
func addPeerRepo(p Peer, t *Torrent)Peer{
	t.Peers= append(t.Peers, p)
	fmt.Println("***addPeerRepo ",t)
	f, err:=os.OpenFile("torrentsFile",os.O_APPEND|os.O_WRONLY,0666)
	if err!=nil {panic (err)}

	out,err:=json.Marshal(t)
	writtenBytes, err := f.WriteString(string(out) )
	fmt.Println("APR wrote %d bytes",writtenBytes)
	f.Close()
	/*fmt.Println("addPeerRepo")
	fmt.Println(t)*/
	return p
}
func GetTorrent(name string) (*Torrent, error) {
	var t Torrent
	torrentF, err := os.Open("torrentsFile")
	if err != nil {
		errors.New("error opening torrentsFile")
	}

	jsonParser := json.NewDecoder(torrentF)
	if err = jsonParser.Decode(&t); err != nil {
		errors.New("parsing config file")
	}
	fmt.Println("--GetTorrent: %v %d %v", t.Name, t.Id, t.Peers)

	torrentF.Close()
	var ret Torrent
	/*for _, torrent := range torrents{
		if strings.Compare(torrent.Name,name) == 0{
			fmt.Println("-Get Torrent ",torrent)
			return &torrent, nil
		}
	}*/
	//return &ret, errors.New("name does not match any torrent")
	return &ret, nil
}
func getIPsRepo(t *Torrent)[]string{
	var ret []string
	for _, peer:= range t.Peers{
		ret = append(ret, peer.IP)
	}
	return ret
}
