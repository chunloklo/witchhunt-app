package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := newHub()
	go hub.run()

	game := newGame(hub)
	go game.run()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, game, c.Writer, c.Request)
	})
	r.Run()
}
