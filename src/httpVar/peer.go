package httpVar

import "sync"
var CurrentPart = 0
var Mutex = &sync.Mutex{}
