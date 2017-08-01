package httpGo
import (
	//"sync"
)
//var Accounts = make(map[string]Account)
//go:generate msgp

type Account struct {
	Name string	`msg:"name"`
	Containers map[string]Container	`msg:"containers"`
	//sync.Mutex
}