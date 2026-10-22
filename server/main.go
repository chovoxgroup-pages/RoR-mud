package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"RoR-mud/world_data"
)

var (
	world    = make(map[string]world_data.Room)
	items    = make(map[string]world_data.Item)
	monsters = make(map[string]world_data.Monster)

	currentWorldState = "default"

	// Tracks living monsters in the world
	activeMobs = make(map[string]*ActiveMob)

	// Upgrades standard HTTP requests to WebSockets
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// ActiveMob represents a living instance of a monster template
type ActiveMob struct {
	InstanceID  string
	Name        string
	Description string
	RoomID      string
	HP          int
	MaxHP       int
}

// ServerPacket represents the dual-channel envelope (Text or GMCP)
type ServerPacket struct {
	Type    string      `json:"type"`    // "text" or "gmcp"
	Module  string      `json:"module"`  // e.g., "Room.Info", "Char.Vitals"
	Payload interface{} `json:"payload"` // HTML string or JSON object
}

// ClientSession holds a single player's WebSocket connection
type ClientSession struct {
	Conn    *websocket.Conn
	Char    *world_data.Player
	WriteMu sync.Mutex
}

// SendText pushes narrative HTML to the main scrolling terminal
func (s *ClientSession) SendText(html string) {
	s.WriteMu.Lock()
	defer s.WriteMu.Unlock()
	s.Conn.WriteJSON(ServerPacket{
		Type:    "text",
		Payload: html,
	})
}

// SendGMCP pushes silent JSON data to update the browser UI panels
func (s *ClientSession) SendGMCP(module string, data interface{}) {
	s.WriteMu.Lock()
	defer s.WriteMu.Unlock()
	s.Conn.WriteJSON(ServerPacket{
		Type:    "gmcp",
		Module:  module,
		Payload: data,
	})
}

func main() {
	loadDatabase()

	// Spawn the initial enemies from monsters.json
	if template, exists := monsters["goblin_warrior_02"]; exists {
		activeMobs["goblin_1"] = &ActiveMob{
			InstanceID:  "goblin_1",
			Name:        "A wide-eyed goblin pyromaniac", // Custom display name
			Description: template.Description,
			RoomID:      "sandpoint_square_01",           // NW Corner
			HP:          template.HP,
			MaxHP:       template.MaxHP,
		}
	}

	if template, exists := monsters["goblin_warrior"]; exists {
		activeMobs["goblin_2"] = &ActiveMob{
			InstanceID:  "goblin_2",
			Name:        "A shrieking goblin warrior", 
			Description: template.Description,
			RoomID:      "sandpoint_square_05",           // Center Square
			HP:          template.HP,
			MaxHP:       template.MaxHP,
		}
	}

	// The 250ms Real-Time Engine Heartbeat
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		for range ticker.C {
			// Combat loops, balance recovery, and AI ticks will go here
		}
	}()

	// Serve the frontend HTML/JS from the /web directory
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Handle incoming WebSocket connections
	http.HandleFunc("/ws", handleConnection)

	log.Println("Runelords MUD Engine listening on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create a dummy character for testing
	char := &world_data.Player{
		ID:           "player_1",
		Name:         "Magus",
		RoomID:       "sandpoint_square_01",
		HP:           24,
		MaxHP:        24,
		Mana:         15, // Spell Points for psionic-style casting
		MaxMana:      15,
		Endurance:    100,
		MaxEndurance: 100,
	}

	session := &ClientSession{Conn: conn, Char: char}
	defer conn.Close()

	session.SendText("<div style='color: #c586c0; font-weight: bold;'>=== RISE OF THE RUNELORDS: BURNT OFFERINGS ===</div>")
	executeLook(session)
	PushVitals(session)

	// Listen for incoming commands from the browser
	for {
		var msg map[string]string
		if err := conn.ReadJSON(&msg); err != nil {
			break // Disconnect
		}

		cmd := strings.TrimSpace(msg["cmd"])
		if cmd != "" {
			// Echo the command back to the terminal
			session.SendText(fmt.Sprintf("<span style='color: #6c7a89;'>&gt; %s</span>", cmd))
			parseCommand(session, cmd)
		}
	}
}

func parseCommand(s *ClientSession, input string) {
	parts := strings.Fields(strings.ToLower(input))
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "l", "look":
		if len(parts) > 1 {
			executeLookTarget(s, parts[1]) // If they typed "look goblin", or "look barrel"
		} else {
		executeLook(s)
		}
	case "n", "north", "s", "south", "e", "east", "w", "west", "ne", "nw", "se", "sw", "u", "up", "d", "down", "in", "out":
		executeMove(s, parts[0])
	default:
		s.SendText("Huh? Type <b>look</b> or use a compass direction.")
	}
}

func executeLook(s *ClientSession) {
	room, exists := world[s.Char.RoomID]
	if !exists {
		s.SendText("You are floating in the void. (Room not found)")
		return
	}

	// Build the text payload
	html := fmt.Sprintf("<div style='color: #5ce6cd; font-weight: bold; font-size: 16px; margin-top: 6px;'>[%s]</div>", room.Name)

	desc := ""
	if stateText, ok := room.Descriptions[currentWorldState]; ok && currentWorldState != "default" {
		desc += "<span style='color: #e74c3c;'>" + stateText + "</span> "
	}
	if defaultText, ok := room.Descriptions["default"]; ok {
		desc += "<span style='color: #d0d4dc;'>" + defaultText + "</span>"
	}
	html += fmt.Sprintf("<div style='margin: 4px 0;'>%s</div>", desc)

	// Inject Active Mobs into the room description
	mobHtml := ""
	for _, mob := range activeMobs {
		if mob.RoomID == s.Char.RoomID {
			mobHtml += fmt.Sprintf("<div style='color: #f1c40f; margin-bottom: 6px;'>%s is here. (HP: %d/%d)</div>", mob.Name, mob.HP, mob.MaxHP)
		}
	}
	html += mobHtml

	// Build Exits
	exitList := []string{}
	for dir := range room.Exits {
		exitList = append(exitList, dir)
	}
	if len(exitList) == 0 {
		exitList = append(exitList, "None")
	}
	html += fmt.Sprintf("<div style='color: #7f8c8d; font-size: 13px;'>Exits: %s</div><br>", strings.Join(exitList, ", "))

	// Send the narrative text to the terminal
	s.SendText(html)

	// Send the silent GMCP payload to update the Automapper UI
	s.SendGMCP("Room.Info", map[string]interface{}{
		"id":    room.ID,
		"name":  room.Name,
		"exits": exitList,
	})
}

func executeMove(s *ClientSession, direction string) {
	dirMap := map[string]string{
		"n": "north", "s": "south", "e": "east", "w": "west",
		"nw": "northwest", "ne": "northeast", "sw": "southwest", "se": "southeast",
		"u": "up", "d": "down", "in": "in", "out": "out",
	}
	shortDirs := map[string]string{
		"north": "n", "south": "s", "east": "e", "west": "w",
		"northwest": "nw", "northeast": "ne", "southwest": "sw", "southeast": "se",
		"up": "u", "down": "d", "in": "in", "out": "out",
	}

	if fullDir, ok := dirMap[direction]; ok {
		direction = fullDir
	}

	room := world[s.Char.RoomID]
	targetRoomID, exists := room.Exits[direction]
	if !exists {
		targetRoomID, exists = room.Exits[shortDirs[direction]]
	}

	if !exists {
		s.SendText("You cannot go that way.")
		return
	}

	// Trap for void falls
	if _, targetExists := world[targetRoomID]; !targetExists {
		s.SendText(fmt.Sprintf("<span style='color: #e74c3c;'>[ENGINE ERROR] The exit points to '%s', but that ID is not loaded.</span>", targetRoomID))
		return
	}

	s.SendText(fmt.Sprintf("You walk %s.", direction))
	s.Char.RoomID = targetRoomID
	executeLook(s)
}

func loadDatabase() {
	worldFile, err := os.ReadFile("world_data/world.json")
	if err != nil {
		log.Printf("Error reading world.json: %v\n", err)
	} else {
		if err = json.Unmarshal(worldFile, &world); err != nil {
			log.Printf("CRITICAL JSON ERROR in world.json: %v\n", err)
		}
	}

	itemFile, err := os.ReadFile("world_data/items.json")
	if err == nil { json.Unmarshal(itemFile, &items) }

	monsterFile, err := os.ReadFile("world_data/monsters.json")
	if err == nil { json.Unmarshal(monsterFile, &monsters) }

	fmt.Printf("Database loaded: %d Rooms | %d Items | %d Monsters\n", len(world), len(items), len(monsters))
}

func PushVitals(s *ClientSession) {
	data := map[string]int{
		"hp":      s.Char.HP,
		"max_hp":  s.Char.MaxHP,
		"mp":      s.Char.Mana,
		"max_mp":  s.Char.MaxMana,
		"end":     s.Char.Endurance,
		"max_end": s.Char.MaxEndurance,
	}
	s.SendGMCP("Char.Vitals", data)
}

func executeLookTarget(s *ClientSession, target string) {
	for _, mob := range activeMobs {
		// Check if the mob is in the same room AND its name contains the target word
		if mob.RoomID == s.Char.RoomID && strings.Contains(strings.ToLower(mob.Name), target) {
			html := fmt.Sprintf("<div style='color: #f1c40f; font-weight: bold;'>%s</div>", mob.Name)
			// Render the specific JSON description
			html += fmt.Sprintf("<div style='color: #d0d4dc; margin: 4px 0;'>%s</div>", mob.Description)
			s.SendText(html)
			return
		}
	}

	s.SendText(fmt.Sprintf("You do not see '%s' here.", target))
}