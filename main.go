package main

import (
	"log"

	"github.com/tjsampson/chat-lacks/core"
	"github.com/tjsampson/chat-lacks/server"
)

func main() {
	logger, file := core.CreateLogger(core.AppConfiguration.LogOutput)
	if file != nil {
		defer file.Close()
	}
	log.Panicln(server.ListenAndServe(logger, core.AppConfiguration))
}
