package httpGo
import (
	"sync"
)
//var Accounts = make(map[string]Account)
var Accounts = struct{
	m map[string]Account
	sync.Mutex
}{m: make(map[string]Account) }

type Account struct {
	Name string
	Containers map[string]Container
	sync.Mutex
}