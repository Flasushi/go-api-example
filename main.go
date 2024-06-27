package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var (
	dataFile = "data.json"
	items    = make(map[string]Item)
	mu       sync.Mutex
)

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mu.Lock()
		defer mu.Unlock()
		json.NewEncoder(w).Encode(items)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		mu.Lock()
		item, ok := items[id]
		mu.Unlock()
		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		var item Item
		json.NewDecoder(r.Body).Decode(&item)
		mu.Lock()
		items[item.ID] = item
		saveItems()
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "item created: %v %v %v", item.ID, item.Name, item.Price)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		mu.Lock()
		delete(items, id)
		saveItems()
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "item deleted : %v", id)
	case http.MethodPut:
		id := r.URL.Query().Get("id")
		var item Item
		json.NewDecoder(r.Body).Decode(&item)
		mu.Lock()
		items[id] = item
		saveItems()
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "item updated: %v", item.ID, item.Name, item.Price)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func loadItems() {
	file, err := os.Open(dataFile)
	if err != nil {
		return
	}
	defer file.Close()
	bytes, _ := os.ReadFile("./db/data.json")
	json.Unmarshal(bytes, &items)
}

func saveItems() {
	bytes, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	os.WriteFile("./db/data.json", bytes, 0644)
}

func main() {
	fmt.Println("sever start....")
	http.HandleFunc("/items", itemsHandler)
	http.ListenAndServe(":8080", nil)
}
