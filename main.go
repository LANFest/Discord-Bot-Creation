package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/commands/admin"
	"github.com/LANFest/Discord-Bot-Creation/commands/chapter"
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	globalData := data.Globals()
	file, readError := ioutil.ReadFile("token.txt")
	utils.Assert("Error reading token file", readError)

	globalData.Token = string(file)

	discord, err := discordgo.New("Bot " + globalData.Token)
	utils.Assert("error creating discord session", err)

	discord.AddHandler(coreMessageHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Insaniquarium!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}

		// Who are we?
		globalData.Bot = discord.State.Ready.User
		fmt.Print("Bot Connected!\n")

		// Who's the owner?
		application, appError := discord.Application("@me")
		utils.Assert("Could not find application!", appError)
		globalData.Owner = application.Owner
		fmt.Printf("Application: %s - Owner: %s\n", application.Name, application.Owner.String())
		fmt.Printf("User: %s -  ID: %s\n-----------------------------------------\n", globalData.Bot.String(), globalData.Bot.ID)

		// Where are we?
		servers := discord.State.Guilds
		fmt.Printf("Servers (%d):\n", len(servers))
		for _, server := range servers {
			fmt.Printf("%s - %s\n", server.Name, server.ID)
		}

	})

	utils.ReadConfig()

	// Set up command handlers
	globalData.CommandHandlers = []interface{}{chapter.PartyOnCommandHandler, admin.WriteConfigCommandHandler, admin.ShutdownCommandHandler}

	err = discord.Open()
	utils.Assert("Error opening connection to Discord", err)
	defer discord.Close()

	<-make(chan struct{})
}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		//Do nothing because a bot is talking
		return
	}

	if !strings.HasPrefix(message.Content, data.Constants().CommandPrefix) {
		// It's not a command, nothing to do here.
		return
	}

	for _, handler := range data.Globals().CommandHandlers {
		handler.(func(*discordgo.Session, *discordgo.MessageCreate))(session, message)
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}
