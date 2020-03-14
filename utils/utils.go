package utils

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Assert : if error exists, log it. Also, conditionally panic.
func Assert(msg string, err error, shouldPanic bool) {
	if err != nil {
		LPrintf("%s: %+v", msg, err)
		if shouldPanic {
			panic(err)
		}
	}
}

// Shutdown : Shuts down the bot
func Shutdown(session *discordgo.Session) {
	LPrint("Shutting Down!")
	session.Logout()
	session.Close()
	os.Exit(0)
}

// LPrint : wrapper for fmt.Print that appends a newline
func LPrint(message string) {
	fmt.Print(message + "\n")
}

// LPrintf : wrapper for fmt.Printf that appends a newline
func LPrintf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

// LogErrorf : LogError but with formatting!
func LogErrorf(componentName string, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	LogError(componentName, message)
}

// LogError : Logs an error
func LogError(componentName string, message string) {
	LPrintf("%s: %s", componentName, message)
}
