package vars


var Foo int
var CurrentId int
var TorrentMap map[string]Torrent

func init() {
	Foo = 5
	CurrentId=0
	TorrentMap=make(map[string]Torrent)

}

func Upd(){
	Foo=Foo+1
}

