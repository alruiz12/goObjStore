package httpVar

import "sync"
var CurrentPart = 0
var P2pPart = 0
var TrackerMutex = &sync.Mutex{}
var PeerMutex = &sync.Mutex{}
var SendMutex = &sync.Mutex{}
var Peer1 string
var Peer2 string
var Peer3 string
var Peers []string
