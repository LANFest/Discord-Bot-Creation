package main

import (
	"io/ioutil"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/admin"
	"github.com/LANFest/Discord-Bot-Creation/chapter"
	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/user"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/ghodss/yaml"
)

var (
	gameLib map[string]string
)

func main() {
	globalData := config.Globals()
	file, readError := ioutil.ReadFile("token.txt")
	utils.Assert("Error reading token file", readError, true)

	globalData.Token = string(file)

	discord, err := discordgo.New("Bot " + globalData.Token)
	//discord.Debug = true
	utils.Assert("error creating discord session", err, true)

	gameLibFile, err := ioutil.ReadFile("games.yml")
	var yamlErr error
	if err != nil {
		yamlErr = yaml.Unmarshal(gameLibFile, &gameLib)
	}
	if err != nil || yamlErr != nil {
		utils.LPrintf("%s %s Error loading game library config file. LFG and game matching features will not work.", err, yamlErr)
	}

	discord.AddHandler(coreReadyHandler)
	discord.AddHandler(coreMessageHandler)
	discord.AddHandler(coreReactionAddHandler)
	discord.AddHandler(coreReactionRemoveHandler)

	utils.ReadConfig()

	// Set up command handlers
	globalData.CommandHandlers = []interface{}{
		chapter.PartyOnCommandHandler,
		admin.WriteConfigCommandHandler,
		admin.ShutdownCommandHandler,
		user.LFGCommandHandler,
	}
	globalData.ReactionAddHandlers = []interface{}{
		user.LFGChannelMessageReactionAdd,
	}

	err = discord.Open()
	utils.Assert("Error opening connection to Discord", err, true)
	defer discord.Close()

	<-make(chan struct{})
}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		//Do nothing because a bot is talking
		return
	}

	if strings.HasPrefix(message.Content, config.Constants().GuildCommandPrefix) && utils.IsGuildMessage(message.Message) {
		// It's a guild command, run through the handlers
		for _, handler := range config.Globals().GuildCommandHandlers {
			// Handlers will return true if they 'handled' the message.
			// This will allow us to circuit-break when we hit the right handler.
			if handler.(func(*discordgo.Session, *discordgo.MessageCreate) bool)(session, message) {
				break
			}
		}
	} else if strings.HasPrefix(message.Content, config.Constants().DMCommandPrefix) && utils.IsDM(message.Message) {
		// It's a DM command, run through the handlers.
		for _, handler := range config.Globals().DMCommandHandlers {
			// Handlers will return true if they 'handled' the message.
			// This allows us to break the circuit once handled.
			if handler.(func(*discordgo.Session, *discordgo.MessageCreate) bool)(session, message) {
				break
			}
		}
	}

	if config.Constants().DebugOutput {
		utils.LPrintf("Message: %+v || From: %s", message.Message, message.Author)
	}
}

func coreReadyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	globalData := config.Globals()
	globalData.Session = discord
	err := discord.UpdateStatus(0, config.Constants().StatusMessage)
	if err != nil {
		utils.LPrintf("Error attempting to set my status: %s", err)
	}

	// Who are we?
	globalData.Bot = discord.State.Ready.User
	utils.LPrint("Bot Connected!")

	// Who's the owner?
	application, appError := discord.Application("@me")
	utils.Assert("Could not find application!", appError, true)
	globalData.Owner = application.Owner
	utils.LPrintf("Application: %s - Owner: %s", application.Name, application.Owner.String())
	utils.LPrintf("User: %s -  ID: %s\n", globalData.Bot.String(), globalData.Bot.ID)
	utils.LPrint("-----------------------------------------")

	// Where are we?
	servers := discord.State.Guilds
	utils.LPrintf("Servers (%d):", len(servers))
	for _, server := range servers {
		utils.LPrintf("%s - %s", server.Name, server.ID)
		validateGuildCoreData(server, config.FindGuildByID(server.ID)) // utils.FindGuildByID has a side-effect of putting the server into the global collection
	}

	globalData.GuildSetups = config.BuildGuildSetupDataList()
	chapter.PromptSetupSteps()

	utils.WriteConfig()
}

func coreReactionAddHandler(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	user, _ := session.User(reaction.UserID)
	if user.Bot {
		return // Ignore when bots add reactions to things.
	}

	for _, handler := range config.Globals().ReactionAddHandlers {
		// Handlers will return true if they 'handled' the message.
		// This will allow us to circuit-break when we hit the right handler.
		if handler.(func(*discordgo.Session, *discordgo.MessageReactionAdd) bool)(session, reaction) {
			break
		}
	}

	if config.Constants().DebugOutput {
		utils.LPrintf("Reaction Add: %+v || From: %s", reaction, user)
	}
}

func coreReactionRemoveHandler(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	user, _ := session.User(reaction.UserID)
	if user.Bot {
		return // Ignore when bots add reactions to things.
	}

	for _, handler := range config.Globals().ReactionDeleteHandlers {
		// Handlers will return true if they 'handled' the message.
		// This will allow us to circuit-break when we hit the right handler.
		if handler.(func(*discordgo.Session, *discordgo.MessageReactionRemove) bool)(session, reaction) {
			break
		}
	}

	if config.Constants().DebugOutput {
		utils.LPrintf("Reaction Delete: %+v || From: %s", reaction, user)
	}
}

func validateGuildCoreData(guild *discordgo.Guild, guildDataModel *config.GuildData) {
	var foundLFG, foundAnnounce, foundAttendee, foundPastAttendee bool // Are the values in the model good?

	// Run through the channels
	for _, channel := range guild.Channels {
		if channel.ID == guildDataModel.LFGCategoryID { // Found our LFGCategory! Still good.
			foundLFG = true
		} else if channel.ID == guildDataModel.AnnounceChannelID { // Found our AnnounceChannel! Still good.
			foundAnnounce = true
			// } else {
			// 	switch channel.Type {
			// 	case discordgo.ChannelTypeGuildCategory:
			// 		if strings.ToLower(channel.Name) == "lfg" && guildDataModel.LFGCategoryID == "" { // We only want to set if it's blank.
			// 			guildDataModel.LFGCategoryID = channel.ID
			// 		}
			// 		break
			// 	case discordgo.ChannelTypeGuildText:
			// 		if strings.ToLower(channel.Name) == "announcements" && guildDataModel.AnnounceChannelID == "" { // We only want to set if it's blank.
			// 			guildDataModel.AnnounceChannelID = channel.ID
			// 		}
			// 	}
		}
	}

	if !foundLFG { // LFGCategory wasn't found.  Maybe it was deleted off the server?
		guildDataModel.LFGCategoryID = ""
	}

	if !foundAnnounce { // AnnounceChannel wasn't found. Maybe it was deleted off the server?
		guildDataModel.AnnounceChannelID = ""
	}

	for _, role := range guild.Roles {
		if role.ID == guildDataModel.AttendeeRoleID { // Found our AttendeeRole! Still good.
			foundAttendee = true
		} else if role.ID == guildDataModel.PastAttendeeRoleID { // Found our PastAttendeeRole! Still good.
			foundPastAttendee = true
			// } else {
			// 	if strings.ToLower(role.Name) == "attendee" && guildDataModel.AttendeeRoleID == "" { // We only want to set if it's blank.
			// 		guildDataModel.AttendeeRoleID = role.ID
			// 	} else if strings.ToLower(role.Name) == "pastattendee" && guildDataModel.PastAttendeeRoleID == "" { // We only want to set if it's blank.
			// 		guildDataModel.PastAttendeeRoleID = role.ID
			// 	}
		}
	}

	if !foundAttendee {
		guildDataModel.AttendeeRoleID = ""
	}

	if !foundPastAttendee {
		guildDataModel.PastAttendeeRoleID = ""
	}
}
