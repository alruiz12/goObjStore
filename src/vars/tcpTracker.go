package vars
import (
	"sync"
)
var FilesMap = SyncFiles{Content: make(map[string]string)}
type SyncFiles struct {
	Content	map[string]string
	Mutex    sync.RWMutex
}

//--------------------------------------------------------------------------------------
/*
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



}
*/