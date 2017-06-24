package httpGo
import(
	"fmt"
	"github.com/alruiz12/simpleBT/src/httpVar"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"strconv"
)

type NodeNum struct {
	Quantity string		`json:"Quantity"`
	Hash string		`json:"Hash"`
}
type jsonKey struct {
	Key string		`json:"Key"`
}


func StartTracker(nodeList []string, proxyAddr []string){
	var nodeAux httpVar.NodeInfo
	for _, node := range nodeList {
		nodeAux.Url=node
		nodeAux.Busy=false
		httpVar.TrackerNodeList = append(httpVar.TrackerNodeList, nodeAux)
	}

	for _, addr := range proxyAddr {
		httpVar.ProxyAddr = append(httpVar.ProxyAddr, addr)
	}

}

/*
GetNodes is called when a GET requests [TrackerURL]/addTorrent.
Sends new json encoded node list back to the sender
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request
 */
func GetNodes(w http.ResponseWriter, r *http.Request){
	var nodeNum NodeNum
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &nodeNum); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("GetNodes: error unprocessable entity: ",err.Error())
			return
		}
		return
	}
	num, err := strconv.Atoi(nodeNum.Quantity)
	if err != nil {
		fmt.Println("GetNodes: error converting string to int response: ",err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("GetNodes: error unprocessable entity: ",err.Error())
			return
		}
		return
	}
	nodeList:=chooseNodes(num)
	httpVar.MapKeyNodes[nodeNum.Hash]=nodeList

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodeList); err != nil {
		fmt.Println("GetNodes: error encoding response: ",err.Error())
	}
}

func chooseNodes(num int)[]string{
	i:=0
	var response []string
	var busies []string
	for i<num{
		if i==len(httpVar.TrackerNodeList){
			fmt.Println("There is not enough free nodes, adding bussies")
			chooseBusyNodes(num,busies, &response)
			return response
		}
		if httpVar.TrackerNodeList[i].Busy==false {
			response = append(response, httpVar.TrackerNodeList[i].Url )
			httpVar.TrackerNodeList[i].Busy=true
		}else{
			busies=append(busies, httpVar.TrackerNodeList[i].Url)
		}
		i++

	}
	return response
}

func chooseBusyNodes(num int, busies []string, response *[]string){
	j:=0 // iterates through bussies
	for len(*response)<num && j<len(busies){
		*response=append(*response, busies[j])
		j++
	}
	if j==len(busies){fmt.Println(" total node number less than tracker asked ")}
		// tracker will receive less than expected


}

func GetNodesForKey(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var key	jsonKey
	if err := json.Unmarshal(body, &key); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("GetNodesForKey: error unprocessable entity: ",err.Error())
			return
		}
		return
	}
	//fmt.Println("GetNodesForKey: KEY IS: ",key.Key)
	if err != nil {
		fmt.Println("GetNodesForKey: error converting string to int response: ",err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("GetNodesForKey: error unprocessable entity: ",err.Error())
			return
		}
		return
	}
	nodeList:=httpVar.MapKeyNodes[key.Key]
	//fmt.Println("GetNodesForKey: About to send : ", nodeList)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodeList); err != nil {
		fmt.Println("GetNodesForKey: error encoding response: ",err.Error())
	}
}
