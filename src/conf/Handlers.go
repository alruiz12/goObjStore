package conf
import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
)

func getPeers(t Torrent){
	//peers:=getPeersRepo(t)
}

/*
addTorrent is called when a POST requests 8080/addTorrent.
Adds a new Torrent to the Tracker file of torrents.
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request.
 */
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
	ret:=addTorrentRepo(t)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
	fmt.Println("... addTorrent FINISHES ...")
}

/*
showTorrents is called when a GET requests 8080/addTorrent.
Sends new json encoded torrent back to the sender
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request
 */
func showTorrents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(torrents); err != nil {
		panic(err)
	}
}

/*
addPeer is called when a POST requests 8080/addPeer.
Adds new peer to given torrent
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request
 */
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

	t,err=GetTorrent(pt.TorrentName)
	if err!=nil {panic(err)}

	auxPeer:= Peer{pt.PeerIP}
	ret:=addPeerRepo(auxPeer,t)

	//Todo if torrent doesn't exist ret:=addTorrentRepo(t)?
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
	fmt.Println("... addPeer FINISHES ...")
}

/*
getIPs is called when a POST requests 8080/getIPs.
Returns the peer's IP addresses, from which the given torrent can be downloaded
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request
 */
func getIPs(w http.ResponseWriter, r *http.Request){
	var torrentName Torrent
	var auxTorrent *Torrent

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
	auxTorrent,err=GetTorrent(torrentName.Name)
	if err!=nil {panic(err)}
	ret:=getIPsRepo(auxTorrent)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}