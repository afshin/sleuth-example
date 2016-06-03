package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/afshin/sleuth-example/types"
	"github.com/gorilla/mux"
	"github.com/ursiform/sleuth"
)

var (
	client *sleuth.Client
	data   = make(map[string][]*types.Comment) // Key is article GUID.
)

func init() {
	var err error
	var raw []byte
	if raw, err = ioutil.ReadFile("data.json"); err != nil {
		panic("Could not read data file: " + err.Error())
	}
	var all []*types.Comment
	if err = json.Unmarshal(raw, &all); err != nil {
		panic("Could not parse: " + err.Error())
	}
	for _, comment := range all {
		data[comment.Article] = append(data[comment.Article], comment)
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	log.Println("GET " + req.URL.String())
	response := new(types.CommentResponse)
	guid := mux.Vars(req)["guid"]
	if comments, ok := data[guid]; ok {
		response.Data = comments
		response.Success = true
		res.WriteHeader(http.StatusOK)
	} else {
		response.Success = false
		response.Message = guid + " not found"
		res.WriteHeader(http.StatusNotFound)
	}
	output, _ := json.Marshal(response)
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/comments/{guid}", handler).Methods("GET")

	// In the real world, the Interface field of the sleuth.Config object
	// should be set so that all services are on the same subnet.
	config := &sleuth.Config{Service: "comment-service", Handler: router}
	client, _ = sleuth.New(config)

	fmt.Println("ready...")
	http.ListenAndServe(":9871", router)
}
