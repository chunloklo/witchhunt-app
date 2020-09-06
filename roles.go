package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func OracleStart(game *Game, oracleUuid uuid.UUID, message gin.H) gin.H {

	nonOracleTown := make([]uuid.UUID, 0)
	for uuid := range game.players {
		player := game.players[uuid]
		if player.Alignment == "TOWN" && uuid != oracleUuid {
			nonOracleTown = append(nonOracleTown, uuid)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nonOracleTown), func(i, j int) { nonOracleTown[i], nonOracleTown[j] = nonOracleTown[j], nonOracleTown[i] })

	message["info"] = nonOracleTown[0]
	return message
}
