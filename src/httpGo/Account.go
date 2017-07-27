package httpGo
import (
	"fmt"
	"sync"
)
//var Accounts = make(map[string]Account)
var Accounts = struct{
	m map[string]Account
	sync.Mutex
}{m: make(map[string]Account) }
func CreateAccount2 (name string){
	Accounts.Lock()
	_, exists := Accounts.m[name]
	if !exists {
		m := make(map[string]Container)
		Accounts.m[name]=Account{Containers:m, Name:name}
		fmt.Println("Account ", name, " created")
	}
	Accounts.Unlock()

}

type Account struct {
	Name string
	Containers map[string]Container
	sync.Mutex
}