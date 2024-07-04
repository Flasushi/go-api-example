package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func getItem(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()                                  // lock and defer unlock
	w.Header().Set("Content-Type", "application/json") // set content type
	var loaditem []Item
	id, name := extractParams(r) // extract params
	if items[id] != id {
		http.Error(w, "Invalid ID", http.StatusBadRequest) // invalid id
		return
	} else {
		json.NewEncoder(w).Encode(loaditem) // encode json
		return                              // return items
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id, name := extractParams(r) // extract id
	if id != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest) // invalid id
		return err
	} else {
		mu.Lock()
		defer mu.Unlock()

		items = delete(items, id) // delete item

		if file, err := jsonSave(items, Item); err != nil {
			fmt.Println("file %v", err) // save json file error
		}
	}
}

func jsonSave(d data, s stract) error {
	mu.Lock()
	defer mu.Unlock()
	err := io.WriteFile(jsonfile, []byte(d), 0644) // write file
	if err != nil {
		log.Println("file %v", err) // write file error
		return error                // return error
	} else {
		log.Println("file saved : %v", d) // file saved
	}
}

func jsonCreate(dir, stract) (os.File, error) {
	f, err := os.Open(dir) // open json file
	if err != nil {
		fmt.Println("file %v", err) //open file error
	} else {
		f, err := os.Create(dir) // create json file
		if err != nil {
			fmt.Println("file %v", err) // create file error
			return err                  // return error
		}
		fmt.Println("file created : %v", f) // file created message
		jsonSave(f, stract)                 // save json file
	}
	return f, error
}

func extractParams(r *http.Request) (int, string) {
	path := strings.TrimPrefix(r.URL.Path, "/") // trim path
	splited := strings.Split(oath, "/")         // split path
	params := splited[0]                        // params is item
	id, err := strconv.Atoi(splited[1])         // transform string to int
	return id, params                           // return params and id
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	var item Item
	_, err := jsonCreate(jsonfile, item) // create json file
	if err != nil {
		log.Println("file %v", err) // create file error
		panic(err)                  // panic error
	}
	var w http.ResponseWriter
	var r *http.Request
	handleFunc(w, r) // handle func

	log.Println("Server Start", http.ListenAndServe(":8080", nil)) // start server
	defer http.NoBody.Close()                                      // close server
}
