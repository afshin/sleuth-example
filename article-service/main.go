package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/afshin/sleuth-example/types"
	"github.com/gorilla/mux"
	"github.com/ursiform/sleuth"
)

const commentsURL = "sleuth://comment-service/comments/%s"

var (
	client *sleuth.Client
	data   = make(map[string]*types.Article) // Key is article GUID.
)

func init() {
	var err error
	var raw []byte
	if raw, err = ioutil.ReadFile("data.json"); err != nil {
		panic("Could not read data file: " + err.Error())

	}
	var all []*types.Article
	if err = json.Unmarshal(raw, &all); err != nil {
		panic("Could not parse: " + err.Error())
	}
	for _, article := range all {
		data[article.GUID] = article
	}
}

func getData(guid string, includeComments bool) (article *types.Article) {
	datum, ok := data[guid]
	if !ok {
		return
	}
	// Data source is immutable, so copy the data.
	article = &types.Article{
		GUID:      datum.GUID,
		Byline:    datum.Byline,
		Headline:  datum.Headline,
		URL:       datum.URL,
		Timestamp: datum.Timestamp}
	if !includeComments {
		return
	}
	url := fmt.Sprintf(commentsURL, guid)
	req, _ := http.NewRequest("GET", url, nil)
	if res, err := client.Do(req); err == nil {
		response := new(types.CommentResponse)
		if err := json.NewDecoder(res.Body).Decode(response); err == nil {
			article.Comments = response.Data
		}
	}
	return
}

func handler(res http.ResponseWriter, req *http.Request) {
	log.Println("GET " + req.URL.String())
	response := new(types.ArticleResponse)
	guid := mux.Vars(req)["guid"]
	include := strings.ToLower(req.URL.Query().Get("includecomments"))
	if article := getData(guid, include == "true"); article != nil {
		response.Data = article
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
	router.HandleFunc("/articles/{guid}", handler).Methods("GET")

	// In the real world, the Interface field of the sleuth.Config object
	// should be set so that all services are on the same subnet.
	client, _ = sleuth.New(&sleuth.Config{})
	client.WaitFor("comment-service")

	fmt.Println("ready...")
	http.ListenAndServe(":9872", router)
}
