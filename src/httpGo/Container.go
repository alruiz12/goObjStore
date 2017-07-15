package httpGo
import (

	"sync"
)
type Container struct {
	Objs map[string]bool
	sync.Mutex
}