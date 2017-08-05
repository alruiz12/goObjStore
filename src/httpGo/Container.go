package httpGo
import (

	//"sync"
)
//go:generate msgp
type Container struct {
	Name string	`msg:"name"`
	Objs map[string]string	`msg:"objs"`
	Policy string	`msg:"policy"`
	//sync.Mutex
}