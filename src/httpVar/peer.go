package httpVar

import "sync"
var CurrentPart = 0
var P2pPart = 0
var TrackerMutex = &sync.Mutex{}
var PeerMutex = &sync.Mutex{}
var SendMutex = &sync.Mutex{}
var DirMutex = &sync.Mutex{}
var GetMutex = &sync.Mutex{}
var WFileMutex = &sync.Mutex{}

var HashMap = make(map[string][]bool)


type NodeInfo struct {
	Url []string
	Busy bool
}
var TrackerNodeList []NodeInfo

var MapKeyNodes = make(map[string][][]string)
var ProxyAddr []string
var SendReady = make(chan int, 200)
var SendP2PReady = make(chan int, 50)

