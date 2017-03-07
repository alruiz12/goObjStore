package vars

/*
Type holding data about a peer
Todo: implement it with ipv4 instead of string
*/
type Peer struct {
	IP	string       `json:"IP"`

}

type Peers []Peer
