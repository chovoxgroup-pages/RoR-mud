package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"RoR-mud/world_data"
)

var (
	worldFile   = "world_data/world.json"
	itemFile    = "world_data/items.json"
	monsterFile = "world_data/monsters.json"
	mutex       sync.Mutex
)

func main() {
	// Serve the visual dashboard
	http.Handle("/", http.FileServer(http.Dir("editor/static")))

	// API endpoint to catch incoming Room data
	http.HandleFunc("/api/save_room", saveRoomHandler)
	http.HandleFunc("/api/save_item", saveItemHandler)
	http.HandleFunc("/api/save_monster", saveMonsterHandler)
	fmt.Println("MUD Editor running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func saveRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newRoom world_data.Room
	if err := json.NewDecoder(r.Body).Decode(&newRoom); err != nil {
		http.Error(w, "Invalid JSON structure", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Initialize the map and attempt to load existing world data
	world := make(map[string]world_data.Room)
	file, err := os.ReadFile(worldFile)
	if err == nil && len(file) > 0 {
		json.Unmarshal(file, &world)
	}

	// Insert or overwrite the room using its ID as the key
	world[newRoom.ID] = newRoom

	// Marshal back to JSON with indentation for readability
	data, _ := json.MarshalIndent(world, "", "  ")

	// Write the updated map back to disk
	os.WriteFile(worldFile, data, 0644)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Room '%s' saved successfully!", newRoom.ID)
}

func saveItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newItem world_data.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid JSON structure", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Load existing Item compendium
	items := make(map[string]world_data.Item)
	file, err := os.ReadFile(itemFile)
	if err == nil && len(file) > 0 {
		json.Unmarshal(file, &items)
	}

	// Insert or overwrite the item
	items[newItem.ID] = newItem

	// Marshal and save
	data, _ := json.MarshalIndent(items, "", "  ")
	os.WriteFile(itemFile, data, 0644)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item '%s' saved to compendium!", newItem.ID)
}

func saveMonsterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newMonster world_data.Monster
	if err := json.NewDecoder(r.Body).Decode(&newMonster); err != nil {
		http.Error(w, "Invalid JSON structure", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Load existing Monster compendium
	monsters := make(map[string]world_data.Monster)
	file, err := os.ReadFile(monsterFile)
	if err == nil && len(file) > 0 {
		json.Unmarshal(file, &monsters)
	}

	// Insert or overwrite the monster
	monsters[newMonster.ID] = newMonster

	// Marshal and save
	data, _ := json.MarshalIndent(monsters, "", "  ")
	os.WriteFile(monsterFile, data, 0644)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Monster '%s' saved to compendium!", newMonster.ID)
}
