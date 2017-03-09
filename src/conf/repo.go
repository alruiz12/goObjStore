package conf

import (
	"errors"
	"github.com/alruiz12/simpleBT/src/vars"
	"fmt"
)


/*
addTorrentRepo is called from Handlers.addTorrent after unmarshalling parameter.
Adds a new Torrent to the Tracker file of torrents.
@param1 new torrent to be added
@returns new torrent added
Todo: save and update currentId to file
Todo: identifying torrents by name or id would lead to ensure unique name or id
 */
func addTorrentRepo(t vars.Torrent) vars.Torrent{
	t.Id=vars.CurrentId
	vars.CurrentId++
	vars.TorrentMap[t.Name]=t
	return t
}

/*
addTPeerRepo is called from Handlers.addPeer after unmarshalling parameters.
Adds a new peer to given torrent, saving it back to the Tracker file of torrents.
@param1 new peer to be added
@param2 pointer to the torrent to be added to
return new peer added
Todo: check parameters
*/
func addPeerRepo(p vars.Peer, t *vars.Torrent) (vars.Peer,error){
	t.Peers= append(t.Peers, p)
	vars.TorrentMap[t.Name]=*t
	return p,nil
}

/*
GetTorrent is called from Handlers after unmarshalling parameters.
Searches for a torrent with given name and returns it if found
@param1 name of the torrent
return pointer to the torrent found or error if not found and error
Todo: search by other field (namely ID)
*/
func GetTorrent(name string) (*vars.Torrent, error) {
	var taux vars.Torrent
	var emptyTorrent vars.Torrent
	taux, exists:= vars.TorrentMap[name]
	if !exists {
		fmt.Println("NAME= ",name)
		fmt.Println(vars.TorrentMap)
		return &emptyTorrent, errors.New("name does not match any torrent")
	}
	return &taux, nil
}

/*
getIPsRepo is called from Handlers.getIP after unmarshalling parameters.
Returns a slice of IP addresses, from which the given torrent can be downloaded
@param1 pointer to torrent
return slice of IP addresses
*/
func getIPsRepo(t *vars.Torrent)[]string{
	var ret []string
	for _, peer:= range t.Peers{
		ret = append(ret, peer.IP)
	}
	return ret
}

