package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	jsonfile = "./db/data.json"
	items    = make(map[int]Item)
	mu       sync.Mutex
)

func init() {
	// add item to the map
}

func getItem(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()                                  // lock and defer unlock
	w.Header().Set("Content-Type", "application/json") // set content type
	id, _ := extractParams(r)                          // extract params
	if id != items[id].ID {
		http.Error(w, "Invalid ID", http.StatusBadRequest) // invalid id
	} else {
		json.NewEncoder(w).Encode(items[id]) // encode json
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	id, _ := extractParams(r) // extract id
	if id != items[id].ID {
		http.Error(w, "Invalid ID", http.StatusBadRequest) // invalid id
	}
	delete(items, id)                      // delete item
	jsonSave(items, jsonfile)              // save json
	fmt.Println("Deleted item", items[id]) // deleted item
}

func jsonSave(items map[int]Item, jsonfile string) {
	mu.Lock()
	defer mu.Unlock()
	bytes, _ := json.Marshal(items) // marshal items

	err := os.WriteFile(jsonfile, bytes, 0644) // write file
	if err != nil {
		log.Println("file", err) // write file error
		panic(err)               // panic error
	}
}

func jsonCreate(jsonfile string) *os.File {
	mu.Lock()
	defer mu.Unlock()

	fstat, err := os.Stat(jsonfile) // read file
	if err != nil {
		_, err := os.Create(jsonfile) // create json file
		if err != nil {
			log.Println("file create error", err) // create file error
			panic(err)                            // panic error
		}
	}
	log.Println("file name", os.DirFS(fstat.Name())) // stat file
	file, err := os.Open(jsonfile)
	if err != nil {
		log.Println("aafile", err) // open file error
		panic(err)                 // panic error
	}
	log.Println("file opened", file) // read file
	return file
}

func extractParams(r *http.Request) (int, string) {
	path := strings.TrimPrefix(r.URL.Path, "/") // trim path
	splited := strings.Split(path, "/")         // split path
	id, _ := strconv.Atoi(splited[1])           // transform string to int
	params := splited[2]                        // params is item
	return id, params                           // return params and id
}

func main() {
	f := jsonCreate(jsonfile) // create json file
	defer f.Close()           // defer close file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItem(w, r) // get item
		//case http.MethodPost:
		//createItem(w, r) // create item
		//case http.MethodPut:
		//updateItem(w, r) // update item item = append(items[id, items]) // append items put
		case http.MethodDelete:
			deleteItem(w, r) // delete item
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // method not allowed
		}
	})
	log.Println("Server Start", http.ListenAndServe(":8080", nil)) // start server
}
