package conf

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"bytes"
	"strconv"
	"bufio"
	"strings"
	"errors"
)

var currentId int
var torrents Torrents

/*
addTorrentRepo is called from Handlers.addTorrent after unmarshalling parameter.
Adds a new Torrent to the Tracker file of torrents.
@param1 new torrent to be added
@returns new torrent added
Todo: save and update currentId to file
Todo: identifying torrents by name or id would lead to ensure unique name or id
 */
func addTorrentRepo(t Torrent) Torrent{
	var err error
	currentId += 1
	t.Id = currentId
	torrents = append(torrents, t)

	f, err:=os.OpenFile("torrentsFile",os.O_APPEND|os.O_WRONLY,0666)
	if err!=nil {panic (err)}
	out,err:=json.Marshal(t)
	var buffer bytes.Buffer
	buffer.WriteString(string(out))
	buffer.WriteString(";")
	_, err = f.Write(buffer.Bytes())
	f.Close()

	aux,err:=ioutil.ReadFile("nTorrents")
	if err!=nil {panic (err)}
	nTorrents,err:=strconv.Atoi(string(aux))
	nTorrents++
	if err!=nil {panic (err)}
	output:=bytes.Replace(aux, aux, []byte(  string(nTorrents)), -1)
	if err = ioutil.WriteFile("nTorrents", output, 0666); err != nil{
		panic(err)
	}

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
func addPeerRepo(p Peer, t *Torrent)Peer{
	t.Peers= append(t.Peers, p)
	f, err:=os.OpenFile("torrentsFile",os.O_APPEND|os.O_WRONLY,0666)
	if err!=nil {panic (err)}
	out,err:=json.Marshal(t)
	var buffer bytes.Buffer
	buffer.WriteString(string(out))
	buffer.WriteString(";")
	_, err = f.Write(buffer.Bytes())
	f.Close()
	return p
}

/*
GetTorrent is called from Handlers after unmarshalling parameters.
Searches for a torrent with given name and returns it if found
@param1 name of the torrent
return pointer to the torrent found or error if not found
Todo: search by other field (namely ID)
*/
func GetTorrent(name string) (*Torrent, error) {
	var taux Torrent
	f, err :=os.Open("torrentsFile")
	r:=bufio.NewReader(f)
	line, err:= r.ReadString(';')
	for line!="" {
		cleanLine:=line
		cleanLine=cleanLine[:len(cleanLine)-1]
		err = json.Unmarshal([]byte(cleanLine), &taux);
		if strings.Compare(taux.Name,name)==0{
			output:=bytes.Replace([]byte(line), []byte(line), []byte(  ""), -1)
			if err = ioutil.WriteFile("torrentsFile", output, 0666); err != nil{
				panic(err)
			}
			return &taux, nil
		}
		if err!=nil {panic (err)}
		line, err= r.ReadString('\n')
	}
	var emptyTorrent Torrent
	return &emptyTorrent, errors.New("name does not match any torrent")
}

/*
getIPsRepo is called from Handlers.getIP after unmarshalling parameters.
Returns a slice of IP addresses, from which the given torrent can be downloaded
@param1 pointer to torrent
return slice of IP addresses
*/
func getIPsRepo(t *Torrent)[]string{
	var ret []string
	for _, peer:= range t.Peers{
		ret = append(ret, peer.IP)
	}
	return ret
}


func getPeersRepo(t Torrent) Peers{
	return t.Peers
}

