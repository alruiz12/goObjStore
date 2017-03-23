package conf

import "testing"

func TestAnnounce(t *testing.T){
	StartAnnouncing(2,9)
	CheckInactivePeers(5)
}
