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
func getIPsRepo(t *Torrent)[]string{
	var ret []string
	for _, peer:= range t.Peers{
		ret = append(ret, peer.IP)
	}
	return ret
}
