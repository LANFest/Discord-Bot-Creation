package main

import (
	"io/ioutil"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/admin"
	"github.com/LANFest/Discord-Bot-Creation/chapter"
	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/user"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq"
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
	// Uncomment the below line to have discordgo dump War and Peace into your buffer on every command.
	//discord.Debug = true
	utils.Assert("Error creating discord session", err, true)

	gameLibFile, err := ioutil.ReadFile("games.yml")
	var yamlErr error
	if err != nil {
		yamlErr = yaml.Unmarshal(gameLibFile, &gameLib)
	}
	if err != nil || yamlErr != nil {
		utils.LogErrorf("Main", "%s %s Error loading game library config file. LFG and game matching features will not work.", err, yamlErr)
	}

	discord.AddHandler(coreReadyHandler)
	discord.AddHandler(coreBotJoinHandler)
	discord.AddHandler(coreMessageHandler)
	discord.AddHandler(coreReactionAddHandler)
	discord.AddHandler(coreReactionRemoveHandler)

	config.ReadConfig()

	// Set up command handlers
	globalData.GuildCommandHandlers = []interface{}{
		chapter.PartyOnCommandHandler,
	}

	globalData.DMCommandHandlers = []interface{}{
		admin.WriteConfigDMCommandHandler,
		admin.ShutdownDMCommandHandler,
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

func coreBotJoinHandler(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	utils.LPrintf("Joining %s - %s", guildCreate.Guild.Name, guildCreate.Guild.ID)
	guildData := config.FindGuildByID(guildCreate.Guild.ID)
	validateGuildCoreData(guildCreate.Guild, guildData) // config.FindGuildByID has a side-effect of putting the server into the global collection

	// should we prompt for a config?
	setupStep := config.GetNextGuildSetupStep(guildData)

	if setupStep != config.GuildSetupStepComplete {
		// First, do we already have an existing GuildSetup for this? (stale data, maybe?)
		ownerSetup, ok := linq.From(config.Globals().OwnerSetups).FirstWithT(func(os config.OwnerSetups) bool {
			anyGS := linq.From(os.GuildSetups).AnyWithT(func(gs *config.GuildSetupData) bool { return gs.GuildID == guildCreate.Guild.ID })
			return anyGS
		}).(config.OwnerSetups)

		if !ok {
			// Nothing known, and this is a new server.  Pop them into the list.
			myGuildSetup := new(config.GuildSetupData)
			myGuildSetup.GuildID = guildCreate.Guild.ID
			myGuildSetup.SetupStep = setupStep
			config.UpsertGuildSetup(*myGuildSetup, guildCreate.Guild.OwnerID)
		}

		// We've potentially done an upsert.  We should fetch out of the globals again.
		ownerSetup, ok = linq.From(config.Globals().OwnerSetups).FirstWithT(func(os config.OwnerSetups) bool {
			anyGS := linq.From(os.GuildSetups).AnyWithT(func(gs config.GuildSetupData) bool { return gs.GuildID == guildCreate.Guild.ID })
			utils.LPrintf("%v", anyGS)
			return anyGS
		}).(config.OwnerSetups)

		// It's the first on the list, so prompt again!
		if ownerSetup.GuildSetups[0].GuildID == guildCreate.Guild.ID {
			owner, ownerError := session.User(guildCreate.Guild.OwnerID)
			if ownerError != nil {
				utils.LogErrorf("coreBotJoinHandler", "Unable to find Owner by ID - %s", guildCreate.Guild.OwnerID)
				return
			}
			chapter.PromptSetupStepByUser(owner, guildCreate.Guild, setupStep)
		}
	}
}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		//Do nothing because a bot is talking
		return
	}

	if strings.HasPrefix(message.Content, config.Constants().GuildCommandPrefix) && utils.IsGuildMessage(session, message.Message) {
		// It's a guild command, run through the handlers
		for _, handler := range config.Globals().GuildCommandHandlers {
			// Handlers will return true if they 'handled' the message.
			// This will allow us to circuit-break when we hit the right handler.
			if handler.(func(*discordgo.Session, *discordgo.MessageCreate) bool)(session, message) {
				break
			}
		}
	} else if utils.IsDM(session, message.Message) {
		if strings.HasPrefix(message.Content, config.Constants().DMCommandPrefix) {
			// It's a DM command, run through the handlers.
			for _, handler := range config.Globals().DMCommandHandlers {
				// Handlers will return true if they 'handled' the message.
				// This allows us to break the circuit once handled.
				if handler.(func(*discordgo.Session, *discordgo.MessageCreate) bool)(session, message) {
					break
				}
			}
		} else if linq.From(config.Globals().OwnerSetups).AnyWithT(func(os config.OwnerSetups) bool { return os.OwnerID == message.Author.ID }) {
			// It is potentially a config response.  Run through the handler.
			chapter.ConfigResponseDMHandler(session, message.Author, message)
		}
	}

	if config.Constants().DebugOutput {
		utils.LogErrorf("Main", "Message: %+v || From: %s", message.Message, message.Author)
	}
}

func coreReadyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	globalData := config.Globals()
	globalData.Session = discord
	err := discord.UpdateStatus(0, config.Constants().StatusMessage)
	utils.Assert("Error attempting to set Status", err, false)

	// Who are we?
	globalData.Bot = discord.State.Ready.User
	utils.LogError("Main", "Bot Connected!")

	// Who's the owner?
	application, appError := discord.Application("@me")
	utils.Assert("Could not find application!", appError, true)
	globalData.Owner = application.Owner
	utils.LPrintf("Application: %s - Owner: %s", application.Name, application.Owner.String())
	utils.LPrintf("User: %s -  ID: %s\n", globalData.Bot.String(), globalData.Bot.ID)
	utils.LPrintf("-----------------------------------------")

	// Where are we?
	servers := discord.State.Guilds
	utils.LPrintf("Servers (%d):", len(servers))
	for _, server := range servers {
		utils.LPrintf("%s - %s", server.Name, server.ID)
		validateGuildCoreData(server, config.FindGuildByID(server.ID)) // config.FindGuildByID has a side-effect of putting the server into the global collection
	}

	config.BuildOwnerSetupDataList()
	chapter.PromptSetupSteps("")

	config.WriteConfig()
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
		return // Ignore when bots remove reactions from things.
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
	var foundLFG, foundAnnounce, foundAttendee, foundPastAttendee, foundAuthorizedUserID bool // Are the values in the model good?

	// Run through the channels
	for _, channel := range guild.Channels {
		if channel.ID == guildDataModel.LFGCategoryID { // Found our LFGCategory! Still good.
			foundLFG = true
		} else if channel.ID == guildDataModel.AnnounceChannelID { // Found our AnnounceChannel! Still good.
			foundAnnounce = true
		}

		// Found 'em both, stop searching
		if foundLFG && foundAnnounce {
			break
		}
	}

	if !foundLFG { // LFGCategory wasn't found.  Maybe it was deleted off the server?
		guildDataModel.LFGCategoryID = ""
	}

	if !foundAnnounce { // AnnounceChannel wasn't found. Maybe it was deleted off the server?
		guildDataModel.AnnounceChannelID = ""
	}

	for _, role := range guild.Roles {
		// If someone uses the same RoleID for both current and past attendees, we need to check separately.
		if role.ID == guildDataModel.AttendeeRoleID { // Found our AttendeeRole! Still good.
			foundAttendee = true
		}

		if role.ID == guildDataModel.PastAttendeeRoleID { // Found our PastAttendeeRole! Still good.
			foundPastAttendee = true
		}

		// Found 'em both, stop searching
		if foundAttendee && foundPastAttendee {
			break
		}
	}

	if !foundAttendee {
		guildDataModel.AttendeeRoleID = ""
	}

	if !foundPastAttendee {
		guildDataModel.PastAttendeeRoleID = ""
	}

	for _, member := range guild.Members {
		if member.User.ID == guildDataModel.AuthorizedUserID { // Found our AuthorizedUser! Still good.
			foundAuthorizedUserID = true
			break
		}
	}

	if !foundAuthorizedUserID {
		guildDataModel.AuthorizedUserID = ""
	}
}
