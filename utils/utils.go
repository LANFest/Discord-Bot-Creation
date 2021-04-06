package utils

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Assert : if error exists, log it. Also, conditionally panic.
func Assert(msg string, err error, shouldPanic bool) {
	if err != nil {
		log.Printf("%s: %+v", msg, err)
		if shouldPanic {
			panic(err)
		}
	}
}

// Shutdown : Shuts down the bot
func Shutdown(session *discordgo.Session) {
	log.Print("Shutting Down!")
	session.Logout()
	session.Close()
	os.Exit(0)
}
