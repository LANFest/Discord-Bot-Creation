package main //declaring that this is the main package

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
) //using the discord go library

var (
	commandPrefix string
	botID         string
	botname       string
) //setting variables? setting variable type?

func main() { // calling the function command, ????
	discord, err := discordgo.New("Bot NTA1MDg3MDMzOTA4Mzk2MDM3.DvRjqw.2qPeYEkKHtPVBHIOP7Oy6Pkcklk")
	errCheck("error creating discord session", err) // records the error if it occurs
	user, err := discord.User("@me")                //grabbing account information
	errCheck("error retrieving account", err)       //check if error occurred

	botID = user.ID //botID is a variable set to the bots information
	botname = "TestBot"
	discord.AddHandler(commandHandler) // a listener that when it picks up a message create it runs the function
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Insaniquarium!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		fmt.Printf("Test Bot has started on %d servers, Congratulations!", len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	//bullshit lane Functional
	for i := 0; i < len(discord.State.Guilds); i++ {
		currentguild := discord.State.Guilds[i]
		println(currentguild.ID)
		channellist, err := discord.GuildChannels(currentguild.ID)
		errCheck("error retrieving channellist", err)
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
				errCheck("error sending message", err)
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
		errCheck("error creating channel", err)
		message, err := discord.ChannelMessageSend(ownerchannel.ID, "Welcome to the LANFest Discord bot! You are registered as the owner of "+currentguild.Name+", so you will need to answer a few questions to complete setup. If you’d rather not, please type the discord ID of another admin-level access user (e.g. "+botname+botID+") to complete setup for you. Type y to continue.")
		errCheck("error sending message", err)
		println(message.ID)
	}

	defer discord.Close()

	commandPrefix = "!" // this is setting the command prefix which lets the bot know that a command will follow? Starts the command?

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
		//Do nothing because the bot is talking?
		return
	}

	//fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
	//fmt.Printf("\n type of message: %s\n", reflect.TypeOf(message.Message))
	fmt.Printf("\n content of message: %s\n", message.Message.Content)
}
