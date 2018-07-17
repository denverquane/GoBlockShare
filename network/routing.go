package network

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	//"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"github.com/denverquane/GoBlockShare/files"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

var GlobalBlockchain *blockchain.BlockChain
var Torrents []files.TorrentFile

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/", handleGetTorrents).Methods("GET")
	muxRouter.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/layers/{layer}", handleGetLayer).Methods("GET")
	//muxRouter.HandleFunc("/addTransaction", handleWriteTransaction).Methods("POST")
	muxRouter.HandleFunc("/addTorrent", handleReceiveTorrent).Methods("POST")

	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, _ *http.Request) {
	// vars := mux.Vars(r)
	if GlobalBlockchain == nil {
		fmt.Println("Don't have blockchain; making new one")
		temp := blockchain.MakeInitialChain()
		GlobalBlockchain = &temp
	}

	data, err := json.MarshalIndent(*GlobalBlockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("GET chain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleGetLayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	layerId := vars["layer"]
	if Torrents == nil {
		http.Error(w, "Don't have any torrents", http.StatusInternalServerError)
		return
	}

	for _, torr := range Torrents {
		for key, meta := range torr.GetLayerHashMap() {
			if key == layerId {
				file, err := os.Open(torr.GetUrl())
				defer file.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data := make([]byte, meta.Offset)

				file.ReadAt(data, meta.Begin)
				io.WriteString(w, string(data))
			}
		}
	}

	fmt.Println("GET layer: " + layerId)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func handleGetTorrents(w http.ResponseWriter, _ *http.Request) {
	if Torrents == nil {
		fmt.Println("Don't have torrents; making new array")
		Torrents = make([]files.TorrentFile, 0)
	}

	data, err := json.MarshalIndent(Torrents, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("GET torrents")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleReceiveTorrent(w http.ResponseWriter, r *http.Request) {
	var message files.TorrentFile

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	respondWithJSON(w, r, http.StatusCreated, message)
	fmt.Println(message.Validate())
}

/* Below is an example of the input format for writing a transaction via the REST API:

{
"origin":
{
"address":"R9UtQ3QE4NrCxGuriwbI0qWCq0u7WqvjU0Q6muEd9Vk=",
"pubkeyx":86420643971005095497364485743353327828044563134904564182951237567725951244265,
"pubkeyy":84350736413375414420184852907452573247898047974475373171004335402121461174787
},
"txref":[],
"currency": "REP",
"quantity":5.99,
"payload":"Sending!",
"r":67869825206353784434575061723707880946031772528032340694185580017437536660581,
"s":6863529193914569235297749315606845644057909902475373433228108461283191248618,
"destAddr":"R9UtQ3QE4NrCxGuriwbI0qWCq0u7WqvjU0Q6muEd9Vk="
}

*/

//func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
//	// vars := mux.Vars(r)
//
//	var m transaction.RESTWrappedFullTransaction
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&m); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
//		return
//	}
//	defer r.Body.Close()
//
//	trans, _ := m.ConvertToFull()
//	fmt.Println(trans.SignedTrans.ToString())
//	if !transaction.Verify(trans.SignedTrans) {
//		respondWithJSON(w, r, http.StatusBadRequest, "Transaction provided is invalid")
//		return
//	}
//
//	message, success := GlobalBlockchain.AddTransaction(trans, trans.SignedTrans.GetOrigin().Address)
//	if !success {
//		respondWithJSON(w, r, http.StatusBadRequest, message)
//	} else {
//		respondWithJSON(w, r, http.StatusCreated, message)
//	}
//	// BroadcastToAllPeers([]string{"http://localhost:8050/" + vars["channel"] + "/chain"}, newChain)
//}

func respondWithJSON(w http.ResponseWriter, _ *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Write(response)
}
