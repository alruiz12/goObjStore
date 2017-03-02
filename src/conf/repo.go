package conf

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
	"bytes"
	"strconv"
	"bufio"
	"strings"
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

	aux,err:=ioutil.ReadFile("nTorrents")
	if err!=nil {panic (err)}
	nTorrents,err:=strconv.Atoi(string(aux))
	nTorrents++
	if err!=nil {panic (err)}
	output:=bytes.Replace(aux, aux, []byte(  string(nTorrents)), -1)
	if err = ioutil.WriteFile("nTorrents", output, 0666); err != nil{
		panic(err)
	}


	f.Close()
	fmt.Println("		nTorrents= ",nTorrents)
	return t
}
func addPeerRepo(p Peer, t *Torrent)Peer{
	t.Peers= append(t.Peers, p)
	fmt.Println("\n		***addPeerRepo ",t)
	f, err:=os.OpenFile("torrentsFile",os.O_APPEND|os.O_WRONLY,0666)
	if err!=nil {panic (err)}
/*
	aux,err:=ioutil.ReadFile("torrentsFile")
	if err!=nil {panic (err)}
	fmt.Println("read================",string(aux))
	var taux Torrent
	err = json.Unmarshal(aux, &taux);
	if err!=nil {panic (err)}
	fmt.Println(taux.Id)
	fmt.Println(taux.Name)
	fmt.Println(taux.Peers)
	fmt.Println("fin")
	if strings.Compare( taux.Name, t.Name) == 0 {

	}
*/
	out,err:=json.Marshal(t)
	var buffer bytes.Buffer
	buffer.WriteString(string(out))
	buffer.WriteString(";")

	writtenBytes, err := f.Write(buffer.Bytes())
	fmt.Println("		wrote %d bytes",writtenBytes)
	f.Close()
	/*fmt.Println("addPeerRepo")
	fmt.Println(t)*/
	return p
}
func GetTorrent(name string) (*Torrent, error) {
	var taux Torrent
	fmt.Println("		GT====")
	f, err :=os.Open("torrentsFile")
	r:=bufio.NewReader(f)
	line, err:= r.ReadString(';')
	fmt.Println(err)
	//fmt.Println("line0: ",line)
	for line!="" {
		fmt.Println("line: ",line)
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
	/*
	aux,err:=ReadLine()  ioutil.ReadFile("torrentsFile")
	if err!=nil {panic (err)}
	fmt.Println("		GT:read==",string(aux))


	n,err:=ioutil.ReadFile("nTorrents")
	if err!=nil {panic (err)}
	nTorrents,err:=strconv.Atoi(string(n))
	if nTorrents<1{

	}
	fmt.Println("		GT:ntorrents ",nTorrents)




	var taux Torrent
	err = json.Unmarshal(aux, &taux);
	if err!=nil {panic (err)}
	fmt.Println("		GT:in vars",taux.Id," ",taux.Name," ",taux.Peers)
	fmt.Println("		name= ",name)
	if strings.Compare(taux.Name, name) == 0 {
		output:=bytes.Replace(aux, aux, []byte(""), -1)
		if err = ioutil.WriteFile("torrentsFile", output, 0666); err != nil{
			panic(err)
		}
		return &taux, nil
	}



	var ret Torrent
	/*for _, torrent := range torrents{
		if strings.Compare(torrent.Name,name) == 0{
			fmt.Println("-Get Torrent ",torrent)
			return &torrent, nil
		}
	}*/
	//return &ret, errors.New("name does not match any torrent")
	var emptyTorrent Torrent
	return &emptyTorrent, nil
}
func getIPsRepo(t *Torrent)[]string{
	var ret []string
	for _, peer:= range t.Peers{
		ret = append(ret, peer.IP)
	}
	return ret
}
