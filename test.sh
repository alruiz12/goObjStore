curl -H "Content-Type: application/json" -d '{"name":"torrent1"}' http://localhost:8080/addTorrent
curl -H "Content-Type: application/json" -d '{"peerIP":"192.168.1.3","torrentName":"torrent1"}' http://localhost:8080/addPeer
curl -H "Content-Type: application/json" -d '{"peerIP":"bbb","torrentName":"torrent1"}' http://localhost:8080/addPeer
sleep 1
curl -H "Content-Type: application/json" -d '{"name":"torrent1"}' http://localhost:8080/getIPs

