package vars

/*
Type Holding data about torrent
*/
type Announcement struct {
	File	string       `json:"file"`
	IP	string    `json:"IP"`
	Status	string	`json:status`
	// 2=interested, 4=have
	// reference: bitorrent.org/beps/bep_0003.html
}