package httpGo
import (

	//"sync"
)
//go:generate msgp
type Container struct {
	Objs map[string]bool	`msg:"objs"`
	//sync.Mutex
}