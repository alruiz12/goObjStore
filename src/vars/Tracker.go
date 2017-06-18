package vars

import (
	"time"
	"sync"
)

var CurrentId int
var TorrentMap map[string]Torrent
var TrackerPort = ":8080"
var TrackerIP = "127.0.0.1"



type IPCounter struct {
	IPAddr string
	Time   time.Time
	Count  int
	TorrentName string
}

type IPCounterMap map[string]IPCounter

type TrackedTorrents struct {
	IPs map[string]IPCounterMap
	Mutex    sync.RWMutex
}

var TrackedTorrentsMap = TrackedTorrents{IPs: make(map[string]IPCounterMap)}

func init() {
	CurrentId=0
	TorrentMap=make(map[string]Torrent)


}

