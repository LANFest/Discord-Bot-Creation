package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/commands/admin"
	"github.com/LANFest/Discord-Bot-Creation/commands/chapter"
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/ghodss/yaml"
)

var (
	commandPrefix string
	botID         string
	botname       string
	gameLib       map[string]string
)

func main() {
	globalData := data.Globals()
	file, readError := ioutil.ReadFile("token.txt")
	utils.Assert("Error reading token file", readError)

	globalData.Token = string(file)

	discord, err := discordgo.New("Bot " + globalData.Token)
	utils.Assert("error creating discord session", err)

	discord.AddHandler(coreMessageHandler)

	user, err := discord.User("@me")              //grabbing account information
	utils.Assert("error retrieving account", err) //check if error occurred
	botID = user.ID                               //botID is a variable set to the bots information
	botname = "LANFestBot"

	gameLibFile, err := ioutil.ReadFile("games.yml")
	var yamlErr error
	if err != nil {
		yamlErr = yaml.Unmarshal(gameLibFile, &gameLib)
	}
	if err != nil || yamlErr != nil {
		fmt.Printf("%s %s Error loading game library config file. LFG and game matching features will not work.\n", err, yamlErr)
	}

	discord.AddHandler(commandHandler) // a listener that when it picks up a message create it runs the function

	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "with your heart <3")
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
		fmt.Printf("LANFest Bot has started on %d servers.\n", len(servers))
		fmt.Printf("\nServer Details\n")
		for i := 0; i < len(discord.State.Guilds); i++ {
			currentguild := discord.State.Guilds[i]
			fmt.Printf("( ) %s [%s]\n", currentguild.Name, currentguild.ID)
			channellist, err := discord.GuildChannels(currentguild.ID)
			utils.Assert("error retrieving channellist", err)
			for a := 0; a < len(channellist); a++ {
				fmt.Printf(" + %s [%s] (%d)\n", channellist[a].ID, channellist[a].Name, channellist[a].Type)
			}
			for r := 0; r < len(currentguild.Roles); r++ {
				currentrole := currentguild.Roles[r]
				fmt.Printf(" %s [%s]\n", currentrole.Name, currentrole.ID)
			}
			ownerchannel, err := discord.UserChannelCreate(currentguild.OwnerID)
			utils.Assert("error creating channel", err)
			message, err := discord.ChannelMessageSend(ownerchannel.ID, "Welcome to the LANFest Discord bot! You are registered as the owner of "+currentguild.Name+", so you will need to answer a few questions to complete setup. If you’d rather not, please type the discord ID of another admin-level access user (e.g. "+botname+botID+") to complete setup for you. Type y to continue.")
			utils.Assert("error sending message", err)
			println(message.ID)
		}
	})

	utils.ReadConfig()

	// Set up command handlers
	globalData.CommandHandlers = []interface{}{chapter.PartyOnCommandHandler, admin.WriteConfigCommandHandler, admin.ShutdownCommandHandler}

	err = discord.Open()
	utils.Assert("Error opening connection to Discord", err)

	for i := 0; i < len(discord.State.Guilds); i++ {
		currentguild := discord.State.Guilds[i]
		println(currentguild.ID)
		channellist, err := discord.GuildChannels(currentguild.ID)
		utils.Assert("error retrieving channellist", err)
		println(channellist)
		var (
			messagechannel string
		)
		for a := 0; a < len(channellist); a++ {
			if channellist[a].Type == discordgo.ChannelTypeGuildText {
				messagechannel = channellist[a].ID
			}
			if channellist[a].Name == "general" {
				message, err := discord.ChannelMessageSend(channellist[a].ID, "Hello Fish!")
				utils.Assert("error sending message", err)
				println(message.ID)
			}
		}
		for r := 0; r < len(currentguild.Roles); r++ {
			currentrole := currentguild.Roles[r]
			println(currentrole.Name)
			println(currentrole.ID)
			if currentrole.Name == "admin" {
				println(messagechannel)
			}
		}
		ownerchannel, err := discord.UserChannelCreate(currentguild.OwnerID)
		utils.Assert("error creating channel", err)
		message, err := discord.ChannelMessageSend(ownerchannel.ID, "Welcome to the LANFest Discord bot! You are registered as the owner of "+currentguild.Name+", so you will need to answer a few questions to complete setup. If you’d rather not, please type the discord ID of another admin-level access user (e.g. "+botname+botID+") to complete setup for you. Type y to continue.")
		utils.Assert("error sending message", err)
		println(message.ID)
	}

	defer discord.Close()

	commandPrefix = "!"

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
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		return
	}

	for _, handler := range data.Globals().CommandHandlers {
		handler.(func(*discordgo.Session, *discordgo.MessageCreate))(discord, message)
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)

	if strings.HasPrefix(message.Message.Content, "!lfg ") {
		fmt.Printf("LFG Command string found.\n")
		specialCharsRegex, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			discord.MessageReactionAdd(message.Message.ChannelID, message.Message.ID, "632789111958274068")
		}
		gameSuggestion := strings.ToLower(specialCharsRegex.ReplaceAllString(strings.TrimPrefix(message.Message.Content, "!lfg "), ""))
		discord.ChannelMessageSend(message.Message.ChannelID, fmt.Sprintf("%s has been added to the LFG queue for %s for the next 4 hours.", message.Message.Author.Username, gameSuggestion))
	}
}
