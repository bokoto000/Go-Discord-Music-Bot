package main

import (
	"fmt"
	"strings"

	"disco.bot/src/framework"

	"disco.bot/src/cmd"

	"disco.bot/src/database"

	"github.com/bwmarrin/discordgo"
)

var (
	conf       *framework.Config
	CmdHandler *framework.CommandHandler
	Sessions   *framework.SessionManager
	youtube    *framework.Youtube
	botId      string
	PREFIX     string
)

func init() {
	conf = framework.LoadConfig("config.json")
	PREFIX = conf.Prefix

}

var db *database.Db

func main() {
	db = database.SetDb()
	defer db.Close()

	CmdHandler = framework.NewCommandHandler()
	registerCommands()
	Sessions = framework.NewSessionManager()
	youtube = &framework.Youtube{Conf: conf}
	discord, err := discordgo.New(conf.BotToken)
	if err != nil {
		fmt.Println("Error creating discord session,", err)
		return
	}
	if conf.UseSharding {
		discord.ShardID = conf.ShardId
		discord.ShardCount = conf.ShardCount
	}
	usr, err := discord.User("@me")
	if err != nil {
		fmt.Println("Error obtaining account details,", err)
		return
	}
	botId = usr.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		discord.UpdateStatus(0, conf.DefaultStatus)
		guilds := discord.State.Guilds
		fmt.Println("Ready with", len(guilds), "guilds.")
	})
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	fmt.Println("Started")
	<-make(chan struct{})
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botId || user.Bot {
		return
	}
	content := message.Content
	if len(content) <= len(PREFIX) {
		fmt.Println("Len<=Len")
		return
	}
	if content[:len(PREFIX)] != PREFIX {
		fmt.Println("Pref!=preif")
		return
	}
	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}
	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, found := CmdHandler.Get(name)
	if !found {
		return
	}
	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}
	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}
	ctx := framework.NewContext(discord, guild, channel, user, message, conf, CmdHandler, Sessions, youtube, db)
	ctx.Args = args[1:]
	c := *command
	c(*ctx)
}

func registerCommands() {
	// ??? means I haven't dug in
	// TODO: Consistant order?
	CmdHandler.Register("help", cmd.HelpCommand, "Gives you this help message!")
	CmdHandler.Register("join", cmd.JoinCommand, "Join a voice channel !join attic")
	CmdHandler.Register("leave", cmd.LeaveCommand, "Leaves current voice channel")
	CmdHandler.Register("play", cmd.PlayCommand, "Plays whats in the queue")
	CmdHandler.Register("stop", cmd.StopCommand, "Stops the music")
	CmdHandler.Register("add", cmd.AddCommand, "Add a song to the queue !add <youtube-link>")
	CmdHandler.Register("cheers", cmd.CheersCommand, "Add a song to the queue and cheer !add <youtube-link> <cheers>")
	CmdHandler.Register("queue", cmd.QueueCommand, "Print queue???")
	CmdHandler.Register("clear", cmd.ClearCommand, "empty queue???")
	CmdHandler.Register("youtube", cmd.YoutubeCommand, "")
	CmdHandler.Register("keyword", cmd.AddKeywordCommand, "Adds keyowrd")
}
