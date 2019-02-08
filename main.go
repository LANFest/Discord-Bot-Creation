package main //declaring that this is the main package

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
) //using the discord go library

var (
	commandPrefix string
	botID         string
) //setting variables? setting variable type?

func main() { // calling the function command, ????
	discord, err := discordgo.New("Bot NTA1MDg3MDMzOTA4Mzk2MDM3.DvRjqw.2qPeYEkKHtPVBHIOP7Oy6Pkcklk")
	errCheck("error creating discord session", err) // records the error if it occurs
	user, err := discord.User("@me")                //grabbing account information
	errCheck("error retrieving account", err)       //check if error occurred

	botID = user.ID                    //botID is a variable set to the bots information
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
		//Do nothing because the bot is talking
		return
	}

	//fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
	//fmt.Printf("\n type of message: %s\n", reflect.TypeOf(message.Message))
	fmt.Printf("\n content of message: %s\n", message.Message.Content)
}
