package httpVar

import "sync"
var CurrentPart = 0
var P2pPart = 0
var TrackerMutex = &sync.Mutex{}
var PeerMutex = &sync.Mutex{}
var SendMutex = &sync.Mutex{}
var DirMutex = &sync.Mutex{}
var HashMap = make(map[string][]bool)
var Peer1 string
var Peer2 string
var Peer3 string
var Peers []string



type NodeInfo struct {
	Url string
	Busy bool
}
var TrackerNodeList []NodeInfo
var MapKeyNodes = make(map[string][]string)