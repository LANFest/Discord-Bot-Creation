package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botID         string
)

func main() {
	discord, err := discordgo.New("Bot NTA1MDg3MDMzOTA4Mzk2MDM3.DvRjqw.2qPeYEkKHtPVBHIOP7Oy6Pkcklk")
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Insaniquarium!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		fmt.Printf("Test Bot has started on %d servers", len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}
