package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Player struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Alignment string `json:"alignment"`
	Alive     bool   `json:"alive"`
	Ready     bool   `json:"ready"`
}

const (
	ROLE_PRIEST      = "PRIEST"
	ROLE_JUDGE       = "JUDGE"
	ROLE_GRAVEDIGGER = "GRAVEDIGGER"
	ROLE_APPRENTICE  = "APPRENTICE"
	ROLE_SURVIVALIST = "SURVIVALIST"
	ROLE_DOB         = "DOB"
	ROLE_GAMBLER     = "GAMBLER"
	ROLE_FANATIC     = "FANATIC"
	ROLE_ORACLE      = "ORACLE"
	ROLE_WATCHMAN    = "WATCHMAN"
)

var ROLE_PREGAME = map[string]bool{
	ROLE_APPRENTICE: true,
	ROLE_GAMBLER:    true,
	ROLE_ORACLE:     true,
}

var DEFAULT_ROLE_LIST = [...]string{
	ROLE_PRIEST,
	ROLE_JUDGE,
	ROLE_GRAVEDIGGER,
	ROLE_APPRENTICE,
	ROLE_SURVIVALIST,
	ROLE_DOB,
	ROLE_GAMBLER,
	ROLE_FANATIC,
	ROLE_ORACLE,
	ROLE_WATCHMAN,
}

type SetPlayerParameter struct {
	uuid   uuid.UUID
	player Player
}

type GetPlayerParameter struct {
	uuid       uuid.UUID
	playerChan chan Player
}

const (
	LOBBY = iota
)

type PlayerMessage struct {
	uuid    uuid.UUID
	message Message
}

type Game struct {
	hub     *Hub
	players map[uuid.UUID]Player

	roleToUuid map[string]uuid.UUID
	witches    map[uuid.UUID]Player
	numDays    int

	waitingOn map[uuid.UUID]bool
	nextFunc  func()

	// Keep track of votes
	votes map[uuid.UUID]uuid.UUID

	// Register requests from the clients.
	register chan uuid.UUID

	// Unregister requests from clients.
	unregister chan uuid.UUID

	// Ticks game forward
	processMessage chan PlayerMessage
}

func voidFunc() {

}

func newGame(hub *Hub) *Game {
	return &Game{
		hub:     hub,
		players: make(map[uuid.UUID]Player),

		roleToUuid: make(map[string]uuid.UUID),
		witches:    make(map[uuid.UUID]Player),
		numDays:    0,

		waitingOn: make(map[uuid.UUID]bool),
		nextFunc:  voidFunc,

		votes: make(map[uuid.UUID]uuid.UUID),

		register:       make(chan uuid.UUID),
		unregister:     make(chan uuid.UUID),
		processMessage: make(chan PlayerMessage),
	}
}

func (g *Game) run() {
	for {
		select {

		case playerUuid := <-g.register:
			g.players[playerUuid] = Player{
				Name:      "",
				Role:      "",
				Alignment: "",
				Alive:     true,
				Ready:     false,
			}

			message, _ := json.Marshal(gin.H{
				"action":  "SET_NAME",
				"players": g.getAllPlayerInfo(),
			})

			g.hub.sendMessage <- SendMessageParameter{uuid: playerUuid, message: message}

		case playerUuid := <-g.unregister:
			if _, ok := g.players[playerUuid]; ok {
				delete(g.players, playerUuid)
			}

		case clientMessage := <-g.processMessage:
			if clientMessage.message.Type == "SET_NAME" {

				uuid := clientMessage.uuid
				message := clientMessage.message

				player := g.players[uuid]
				player.Name = message.Data[0]
				g.players[uuid] = player

				players := g.getAllPlayerInfo()

				for p := range g.players {
					if g.players[p].Name != "" {

						returnMessage, _ := json.Marshal(gin.H{
							"action":   "LOBBY",
							"players":  players,
							"selfInfo": g.players[p],
						})
						g.hub.sendMessage <- SendMessageParameter{uuid: p, message: returnMessage}
					} else {
						// make everyone without a name to set names
						returnMessage, _ := json.Marshal(gin.H{
							"action":   "SET_NAME",
							"players":  players,
							"selfInfo": g.players[p],
						})
						g.hub.sendMessage <- SendMessageParameter{uuid: p, message: returnMessage}
					}
				}
			}

			if clientMessage.message.Type == "READY" {
				ready := false
				if len(clientMessage.message.Data) > 0 {
					if clientMessage.message.Data[0] == "true" {
						ready = true
					} else {
						ready = false
					}
				}
				player := g.players[clientMessage.uuid]
				player.Ready = ready
				g.players[clientMessage.uuid] = player

				allReady := true

				for p := range g.players {
					if g.players[p].Name != "" && g.players[p].Ready == false {
						allReady = false
					}
				}

				if allReady {
					// Start Game!
					// Assign Roles:
					roles := DEFAULT_ROLE_LIST[:len(g.getAllPlayerInfo())]
					rand.Seed(time.Now().UnixNano())
					rand.Shuffle(len(roles), func(i, j int) { roles[i], roles[j] = roles[j], roles[i] })

					playerNumber := 0
					for uuid := range g.players {
						if g.players[uuid].Name != "" {

							player := g.players[uuid]
							player.Role = roles[playerNumber]
							player.Alignment = "TOWN"
							g.players[uuid] = player
							g.roleToUuid[player.Role] = uuid
							playerNumber += 1
						}
					}

					witchRoles := make([]string, len(roles)-1)

					index := 0
					for _, role := range roles {
						if role != ROLE_PRIEST {
							witchRoles[index] = role
							index += 1
						}
					}
					rand.Seed(time.Now().UnixNano())
					rand.Shuffle(len(witchRoles), func(i, j int) { witchRoles[i], witchRoles[j] = witchRoles[j], witchRoles[i] })
					numWitches := int(math.Floor(math.Sqrt(float64(len(roles)))))

					for i := 0; i < numWitches; i++ {
						uuid := g.roleToUuid[witchRoles[i]]
						player := g.players[uuid]
						player.Alignment = "WITCH"
						g.players[g.roleToUuid[witchRoles[i]]] = player
						g.witches[uuid] = player
					}

					// Prepare for Day 1

					// Clear waiting on map
					g.waitingOn = make(map[uuid.UUID]bool)
					g.nextFunc = g.checkStartDay1

					for uuid := range g.players {
						player := g.players[uuid]
						// check player exists and of pre-game role
						messageJSON := gin.H{
							"action":   "NIGHT",
							"number":   g.numDays,
							"selfInfo": g.players[uuid],
						}
						if player.Alignment == "WITCH" {
							messageJSON["witchInfo"] = g.witches
						}

						if _, ok := ROLE_PREGAME[player.Role]; ok && player.Name != "" {
							g.waitingOn[uuid] = true
							if player.Role == ROLE_APPRENTICE {
								messageJSON["action"] = "APPRENTICE_START"
							}
							if player.Role == ROLE_GAMBLER {
								messageJSON["action"] = "GAMBLER_START"
							}
							if player.Role == ROLE_ORACLE {
								messageJSON["action"] = "ORACLE_START"
								messageJSON = OracleStart(g, uuid, messageJSON)
							}
						}

						messageString, _ := json.Marshal(messageJSON)
						g.hub.sendMessage <- SendMessageParameter{uuid: uuid, message: messageString}
					}
				} else {
					// Update ready info for players
					players := g.getAllPlayerInfo()
					for p := range g.players {
						if g.players[p].Name != "" {

							returnMessage, _ := json.Marshal(gin.H{
								"action":   "LOBBY",
								"players":  players,
								"selfInfo": g.players[p],
							})
							g.hub.sendMessage <- SendMessageParameter{uuid: p, message: returnMessage}
						}
					}
				}
			}

			// stage_person_somethingelse

			if clientMessage.message.Type == "APPRENTICE_START_ROLE_SELECT" {
				message := clientMessage.message
				uuid := clientMessage.uuid
				player := g.players[clientMessage.uuid]

				// Selected Role
				if len(message.Data) == 1 {
					messageJSON := gin.H{
						"action":      "APPRENTICE_START",
						"nightNumber": 0,
						"selfInfo":    g.players[uuid],
					}

					selectedRole := message.Data[0]

					rolePlayer := g.players[g.roleToUuid[selectedRole]]

					messageJSON["selectedRole"] = selectedRole
					messageJSON["rolePlayer"] = rolePlayer

					if player.Alignment == "WITCH" {
						messageJSON["witchInfo"] = g.witches
					}
					stringJSON, _ := json.Marshal(messageJSON)
					g.hub.sendMessage <- SendMessageParameter{uuid: uuid, message: stringJSON}
				} else if len(message.Data) == 2 {
					// Ready!
					if message.Data[1] == "TRUE" {
						g.waitingOn[uuid] = false
						g.nextFunc()
					}

					// messageJSON := gin.H{
					// 	"action":      "APPRENTICE_START",
					// 	"nightNumber": 0,
					// 	"selfInfo":    g.players[uuid],
					// }

					// selectedRole := message.Data[0]
					// rolePlayer := g.players[g.roleToUuid[selectedRole]]

					// messageJSON["selectedRole"] = selectedRole
					// messageJSON["rolePlayer"] = rolePlayer

					// if player.Alignment == "WITCH" {
					// 	messageJSON["witchInfo"] = g.witches
					// }
					// stringJSON, _ := json.Marshal(messageJSON)
					// g.hub.sendMessage <- SendMessageParameter{uuid: uuid, message: stringJSON}
				}
			}

			if clientMessage.message.Type == "DAY_VOTE" {
				message := clientMessage.message
				playerUuid := clientMessage.uuid

				for uuid := range g.players {
					if g.players[uuid].Name == message.Data[0] {
						// register vote
						g.votes[playerUuid] = uuid
					}
				}

				// nameVotes := make(map[string]string)

				// broadcast message
				for uuid := range g.players {
					player := g.players[uuid]
					if g.players[uuid].Name != "" {

						messageJSON := gin.H{
							"action":   "DAY",
							"number":   g.numDays,
							"selfInfo": g.players[uuid],
							"players":  g.getAllPlayerInfo(),
						}

						if player.Alignment == "WITCH" {
							messageJSON["witchInfo"] = g.witches
						}

						stringJSON, _ := json.Marshal(messageJSON)
						g.hub.sendMessage <- SendMessageParameter{uuid: uuid, message: stringJSON}

					}
				}
			}
		}
	}
}

func (g *Game) getAllPlayerInfo() []Player {
	players := make([]Player, 0, len(g.players))
	for p := range g.players {
		if g.players[p].Name != "" {
			players = append(players, g.players[p])
		}
	}
	return players
}

func (g *Game) checkStartDay1() {
	stillWaiting := false
	for uuid := range g.waitingOn {
		if g.waitingOn[uuid] == true {
			stillWaiting = true
		}
	}
	if !stillWaiting {
		fmt.Println("STARTING DAY 1")

		// Clear Votes for Day 1
		g.votes = make(map[uuid.UUID]uuid.UUID)
		g.numDays += 1

		// Update ready info for players
		for uuid := range g.players {
			player := g.players[uuid]
			if g.players[uuid].Name != "" {

				messageJSON := gin.H{
					"action":   "DAY",
					"number":   g.numDays,
					"selfInfo": g.players[uuid],
					"players":  g.getAllPlayerInfo(),
				}

				if player.Alignment == "WITCH" {
					messageJSON["witchInfo"] = g.witches
				}

				stringJSON, _ := json.Marshal(messageJSON)
				g.hub.sendMessage <- SendMessageParameter{uuid: uuid, message: stringJSON}

			}
		}
	}
}

func (g *Game) checkVote() {
}
