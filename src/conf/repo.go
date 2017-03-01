package conf
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
func addPeersRepo(p Peer, t Torrent)Peer{
	t.Peers= append(t.Peers, p)
	return p
}
