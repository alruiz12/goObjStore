package httpVar

import (
	"sync"
)
var CurrentPart = 0
var P2pPart = 0
var TrackerMutex = &sync.Mutex{}
var PeerMutex = &sync.Mutex{}
var SendMutex = &sync.Mutex{}
var DirMutex = &sync.Mutex{}
var GetMutex = &sync.Mutex{}
var WFileMutex = &sync.Mutex{}
var TotalNumMutex = &sync.Mutex{}

var HashMap = make(map[string][]bool)


type NodeInfo struct {
	Url []string
	Busy bool
}
var TrackerNodeList []NodeInfo

var MapKeys = &sync.Mutex{}
var MapCont = &sync.Mutex{}
var MapAcc = &sync.Mutex{}


var MapKeyNodes = make(map[string][][]string)
var MapContNodes = make(map[string][][]string)
var MapAccNodes = make(map[string][][]string)

var ProxyAddr []string
var SendReady = make(chan int, 180)
var SendP2PReady = make(chan int, 20)

//var Accounts = make(map[string]httpGo.Account)
//var AccountMutex