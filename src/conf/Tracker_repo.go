package conf

import (
	"errors"
	"github.com/alruiz12/simpleBT/src/vars"

	"fmt"
	"time"
)


/*
addTorrentRepo is called from Handlers.addTorrent after unmarshalling parameter.
Adds a new Torrent to the Tracker file of torrents.
@param1 new torrent to be added
@returns new torrent added, error if any
Todo: identifying torrents by name or id would lead to ensure unique name or id
 */
func addTorrentRepo(t vars.Torrent) (vars.Torrent, error){
	torrent:=vars.Torrent{}
	_, exists:= vars.TorrentMap[t.Name]
	if exists {
		return torrent, errors.New("torrent already exists")
	}
	t.Id=vars.CurrentId
	vars.CurrentId++
	vars.TorrentMap[t.Name]=t
	return t, nil
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


/*
TrackPeers is called from Handlers.listenAnnounce after unmarshalling parameters.
@param1 peer IP
@param2 torrent name
*/
func TrackPeers(IP string, torrent string){
	fmt.Println("	...TrackPeers starts...")

	Incr(IP,torrent)

	fmt.Println("	...TrackPeers finishes...")
}


// Increment IP address counter and update time
func Incr(ipAddr string, torrent string) {
	now := time.Now().UTC()
	vars.TrackedTorrentsMap.Mutex.Lock()
	defer 	vars.TrackedTorrentsMap.Mutex.Unlock()

	IPmap, torrentFound := vars.TrackedTorrentsMap.IPs[torrent]
	if !torrentFound { //empty value: inner map
		NewCounter := new(vars.IPCounter)
		NewCounter.IPAddr=ipAddr
		NewCounter.Time=now
		NewCounter.TorrentName=torrent
		IPmap=make(map[string]vars.IPCounter)
		IPmap[ipAddr]=*NewCounter


	} else{ //already a map for given torrent, search given IP and update
		counter, IPexists:= IPmap[ipAddr]
		counter.Time=now
		counter.Count++
		if !IPexists{ //IP not registered
			counter.IPAddr=ipAddr
			counter.TorrentName=torrent

		}
		IPmap[ipAddr]=counter
	}
	vars.TrackedTorrentsMap.IPs[torrent]=IPmap
}


// Delete IP address counter
func Delete(torrentName string,ipAddr string) {
	fmt.Println("Deleting "+ipAddr)
	vars.TrackedTorrentsMap.Mutex.Lock()
	defer 	vars.TrackedTorrentsMap.Mutex.Unlock()

	delete(vars.TrackedTorrentsMap.IPs[torrentName], ipAddr)
}

// Get old IP address counters old durations ago
func OldIPCounters(old time.Duration){
	oldTime := time.Now().UTC().Add(-old)
	vars.TrackedTorrentsMap.Mutex.Lock()
	defer 	vars.TrackedTorrentsMap.Mutex.Unlock()
	for _, torrents := range vars.TrackedTorrentsMap.IPs {
		for _, counter := range torrents{
			if counter.Time.Before(oldTime){
				Delete(counter.TorrentName,counter.IPAddr)
			}
		}

	}
}

func CheckInactivePeers (interval time.Duration) {
	fmt.Println("checking inactive Peers")
	ticker := time.NewTicker(interval * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				OldIPCounters(interval * time.Second)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}