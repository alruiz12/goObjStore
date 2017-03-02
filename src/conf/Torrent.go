package conf


/*
Type Holding data about torrent
Todo: identifying data by unique ID
 */
type Torrent struct {
	Id	int       `json:"id"`
	Name	string    `json:"name"`
	Size 	int      `json:"size"`
	Trackers []int	 `json:"trackers"`
	Peers 	Peers	 `json:"peers"`
}

type Torrents []Torrent