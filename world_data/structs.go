package world_data

// Room represents a physical location within the MUD
type Room struct {
	ID   string `json:"id"`   // Unique identifier for the room
	Name string `json:"name"` // Name of the room
	// Descriptions maps a World State to a specific group of readaloud text
	Descriptions map[string]string `json:"descriptions"` // Text that is read aloud when a player enters the room
	Exits        map[string]string `json:"exits"`        // Maps a direction (e.g., "north") to a Room ID
	NPCs         []string          `json:"npcs"`         // Slice of NPC IDs currenrtly in the room
	Items        map[string]int    `json:"items"`        // Map of Item IDs to quantities "on" the room floor
}

// AbilityScores represents the standard Pathfinder 1e ability scores
type AbilityScores struct {
	STR int `json:"str"` // Strength
	DEX int `json:"dex"` // Dexterity
	CON int `json:"con"` // Constitution
	INT int `json:"int"` // Intelligence
	WIS int `json:"wis"` // Wisdom
	CHA int `json:"cha"` // Charisma
}

// Saves represents standard Pathfinder 1e saving throws
type Saves struct {
	Fortitude int `json:"fortitude"` // Fortitude save
	Reflex    int `json:"reflex"`    // Reflex save
	Will      int `json:"will"`      // Will save
}

// Item represents an object or place of loot in the world
type Item struct {
	ID          string  `json:"id"`          // Unique identifier for the item
	Name        string  `json:"name"`        // Name of the item
	Description string  `json:"description"` // Description of the item
	Weight      float64 `json:"weight"`      // Weight of the item in pounds
	Value       int     `json:"value"`       // Monetary value of the item in gold pieces
}

// ArmorClass represents the split AC mechanics of Pathfinder 1e
type ArmorClass struct {
	Total      int `json:"total"`
	Touch      int `json:"touch"`
	FlatFooted int `json:"flat_footed"`
}

// Monster represents an NPC or creature
type Monster struct {
	ID               string         `json:"id"`   // Unique identifier for the monster
	Name             string         `json:"name"` // Name of the monster
	Description      string         `json:"description"`
	HP               int            `json:"hp"`
	MaxHP            int            `json:"max_hp"`
	Mana             int            `json:"mana"`
	MaxMana          int            `json:"max_mana"`
	Endurance		 int            `json:"endurance"`
	MaxEndurance     int            `json:max_endurance"`
	AC               ArmorClass     `json:"ac"`
	Speed            int            `json:"speed"`
	Senses           []string       `json:"senses"` // e.g., ["darkvision 60', "tremorsense 30"]
	MeleeAttacks     []string       `json:"melee_attacks"`
	RangedAttacks    []string       `json:"ranged_attacks"`
	Saves            Saves          `json:"saves"`    // Saving throws
	Defenses         []string       `json:"defenses"` // e.g., ["immunity to fire", "resistance to cold 10"]
	Abilities        AbilityScores  `json:"abilities"`
	SpecialQualities []string       `json:"special_qualities"` // e.g., ["regeneration 5", "scent"]
	SpecialAttacks   []string       `json:"special_attacks"`   // e.g., ["breath weapon", "spell-like abilities"]
	Inventory        map[string]int `json:"inventory"`         // Maps Item IDs to quantities
}

// Player represents a connected player's character.
type Player struct {
	ID               string         `json:"id"`      // Unique identifier for the player
	Name             string         `json:"name"`    // Name of the player character
	RoomID           string         `json:"room_id"` // ID of the room the player is currently in
	HP               int            `json:"hp"`      // Current hit points
	MaxHP            int            `json:"max_hp"`  // Maximum hit points
	Mana             int            `json:"mana"`
	MaxMana          int            `json:"max_mana"`
	Endurance		 int            `json:"endurance"`
	MaxEndurance     int            `json:max_endurance"`
	AC               ArmorClass     `json:"ac"`      // Armor Class
	Speed            int            `json:"speed"`   // Movement speed
	Senses           []string       `json:"senses"`  // e.g., ["darkvision 60'", "tremorsense 30"]
	MeleeAttacks     []string       `json:"melee_attacks"`
	RangedAttacks    []string       `json:"ranged_attacks"`
	Saves            Saves          `json:"saves"`    // Saving throws
	Defenses         []string       `json:"defenses"` // e.g., ["immunity to fire", "resistance to cold 10"]
	Abilities        AbilityScores  `json:"abilities"`
	SpecialQualities []string       `json:"special_qualities"` // e.g., ["regeneration 5", "scent"]
	SpecialAttacks   []string       `json:"special_attacks"`   // e.g., ["breath weapon", "spell-like abilities"]
	Inventory        map[string]int `json:"inventory"`         // Maps Item IDs to quantities
}
