package conf
import (
	"encoding/json"
	"fmt"
	"net/http"
	//"simpleBT/src/mux"
	"io/ioutil"
	"io"
)

func getPeers(t Torrent){
	//peers:=getPeersRepo(t)
	//fmt.Fprintln(peers)
}
func addTorrent(w http.ResponseWriter, r *http.Request){
	fmt.Println("... addTorrent STARTS ...")
	var t Torrent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &t); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	fmt.Println(t)
	ret:=addTorrentRepo(t)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)//Todo v maybe torrents???????
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
	fmt.Println("... addTorrent FINISHES ...")
}

func showTorrents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(torrents); err != nil {
		panic(err)
	}
}

func addPeer(w http.ResponseWriter, r *http.Request){
	fmt.Println("... addPeer STARTS ...")
	var t *Torrent
	type PeerAndTorrent struct {
		PeerIP string		`json:"peerIP"`
		TorrentName string	`json:"torrentName"`
	}
	var pt PeerAndTorrent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &pt); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
//	p=pt.peer
//	t=pt.torrent
	//fmt.Println(pt.TorrentName)
	t,err=GetTorrent(pt.TorrentName)
	if err!=nil {panic(err)}
	auxPeer:= Peer{pt.PeerIP}
	ret:=addPeerRepo(auxPeer,t)
	fmt.Println("-after addPeerRepo returns ", t)
	fmt.Println(auxPeer)
	//Todo if torrent doesn't exist ret:=addTorrentRepo(t)?
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)//Todo v maybe torrents???????
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
	fmt.Println("... addPeer FINISHES ...")


}

func getIPs(w http.ResponseWriter, r *http.Request){
	var torrentName Torrent
	var auxTorrent *Torrent
	//var auxTorrent Torrent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &torrentName); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	//fmt.Println(torrentName)
	auxTorrent,err=GetTorrent(torrentName.Name)
	fmt.Println(auxTorrent)
	if err!=nil {panic(err)}
	ret:=getIPsRepo(auxTorrent)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}