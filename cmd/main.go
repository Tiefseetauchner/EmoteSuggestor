package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/Tiefseetauchner/EmoteSuggestor/pkg/LocalFunctions"
	"github.com/bwmarrin/discordgo"
)

var commands = []string{"a", "b"}
var commandRegex = regexp.MustCompile(`^[!](\p{L}+)[ ]?(.*)$`)
var helpMsg = "**Help**\n" +
	"```\n" +
	"!suggest [link]: adds an emoji from a link or the attachment\n" +
	"!help: shows this help\n" +
	"```\n"

func main() {
	var config = localFunctions.LoadConfiguration("./config.json")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if []rune(message.Content)[0] == '!' {
		if commandRegex.Match([]byte(message.Content)) {
			var command = commandRegex.FindStringSubmatch(message.Content)[1]
			var arguments = commandRegex.FindStringSubmatch(message.Content)[2]

			runCommand(command, arguments, session, message)
		}
	}
}

func runCommand(command string, arguments string, session *discordgo.Session, message *discordgo.MessageCreate) {
	switch command {
	case "help":
		_, _ = session.ChannelMessageSend(message.ChannelID, helpMsg)
		break
	case "suggest":
		if len(message.Attachments) != 0 {
			_, _ = session.ChannelMessageSend(message.ChannelID, "Suggesting adding: "+message.Attachments[0].URL)
		} else if len(arguments) > 0 {
			args := strings.Split(arguments, " ")
			_, _ = session.ChannelMessageSend(message.ChannelID, "Suggesting adding: "+args[0])
		} else {
			_, _ = session.ChannelMessageSend(message.ChannelID, "No link or Embed provided. Usage:\n"+
				"```\n"+
				"!suggest name [link] ...: suggests adding emote from link with name name to emojis\n"+
				"!suggest name: suggests adding emote from attachment with name name to emojis"+
				"```")
		}
	default:
		_, _ = session.ChannelMessageSend(message.ChannelID, "Command "+command+
			" not found. Type !help for a list of commands")
	}
}
