package conf

import (
	"strings"
	"errors"
	"fmt"
)
//import "fmt"

var currentId int

var torrents Torrents

func getPeersRepo(t Torrent) Peers{
	return t.Peers
}

func addTorrentRepo(t Torrent) Torrent{
	currentId += 1
	t.Id = currentId
	torrents = append(torrents, t)
	return t
}
func addPeerRepo(p Peer, t *Torrent)Peer{
	t.Peers= append(t.Peers, p)
	/*fmt.Println("addPeerRepo")
	fmt.Println(t)*/
	return p
}
func GetTorrent(name string) (*Torrent, error) {
	var ret Torrent
	for _, torrent := range torrents{
		if strings.Compare(torrent.Name,name) == 0{
			fmt.Println("-Get Torrent ",torrent)
			return &torrent, nil
		}
	}
	return &ret, errors.New("name does not match any torrent")
}
func getIPsRepo(t *Torrent)[]string{
	var ret []string
	for _, peer:= range t.Peers{
		ret = append(ret, peer.IP)
	}
	return ret
}
