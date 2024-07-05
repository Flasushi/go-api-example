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
	mu.Lock()                    // lock
	defer mu.Unlock()            // defer unlock
	readItems(jsonfile)          // read items
	id, _ := extractParams(w, r) // extract params
	if id != int(items[id].ID) {
		log.Println("Invalid ID", id, items[id])
		http.Error(w, "Invalid ID", http.StatusBadRequest) // invalid id
	} else {
		json.NewEncoder(w).Encode(items[id]) // encode json
		log.Println("Get item", items[id])   // get item
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	id, _ := extractParams(w, r) // extract id
	if id != items[id].ID {
		http.Error(w, "Invalid ID", http.StatusBadRequest) // invalid id
	}
	delete(items, id)                      // delete item
	jsonSave(items, jsonfile)              // save json
	fmt.Println("Deleted item", items[id]) // deleted item
}

func postItem(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	id, name := extractParams(w, r)
	items := readItems(jsonfile)
	var newItem Item

	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest) // invalid input
		return
	}

	items[newItem.ID] = Item{ID: id, Name: name} // new item

	w.Header().Set("Content-Type", "application/json") // set content type
	jsonSave(items, jsonfile)

	fmt.Println("Posted item", items[id]) // posted item
	json.NewEncoder(w).Encode(items[id])  // encode json
}

func jsonSave(items map[int]Item, jsonfile string) {
	mu.Lock()
	defer mu.Unlock()
	bytes, _ := json.Marshal(items)            // marshal items
	file, _ := os.Open(jsonfile)               // open file
	defer file.Close()                         // close file
	err := os.WriteFile(jsonfile, bytes, 0644) // write file
	if err != nil {
		log.Println("file", err) // write file error
		panic(err)               // panic error
	} else {
		log.Println("file saved", items) // file saved
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
		log.Println("file", err) // open file error
		panic(err)               // panic error
	}
	defer file.Close()
	log.Println("file opened", file) // read file
	return file
}

func extractParams(w http.ResponseWriter, r *http.Request) (int, string) {
	path := strings.TrimPrefix(r.URL.Path, "/") // trim path
	splited := strings.Split(path, "/")         // split path
	fmt.Println(splited)                        // print splited
	id, err := strconv.Atoi(splited[1])         // transform string to int
	name := splited[2]                          // name is item

	if err != nil {
		fmt.Println("hugahuga", r.URL.Path, err, splited[1]) // print error
		http.Error(w, "Invalid ID", http.StatusBadRequest)   // invalid id

	}
	return id, name // retrun id, name
}

func readItems(jsonfile string) map[int]Item {
	file, err := os.Open(jsonfile)
	if err != nil {
		log.Println("open err", err)
	}
	defer file.Close()

	bytedata, err := io.ReadAll(file)
	if err != nil {
		log.Println("read err", err)
	}

	var tempItem []Item
	if err := json.Unmarshal(bytedata, &tempItem); err != nil {
		log.Println("json unmarshal err", err)
	}
	for _, item := range tempItem {
		items[item.ID] = item
	}
	return items
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItem(w, r) // get item
	case http.MethodPost:
		postItem(w, r) // post item
	//case http.MethodPut:
	//updateItem(w, r) // update item item = append(items[id, items]) // append items put
	case http.MethodDelete:
		deleteItem(w, r) // delete item
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // method not allowed
	}
}

func main() {
	jsonCreate(jsonfile) // create json file
	readItems(jsonfile)  // read items

	http.HandleFunc("/GET/", itemHandler)    // get item
	http.HandleFunc("/POST/", itemHandler)   // post item
	http.HandleFunc("/DELETE/", itemHandler) // delete item

	http.ListenAndServe(":8080", nil) // listen and serve
}
