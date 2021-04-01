package main

import (
	"fmt"
	"go/types"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
	"syscall"

	"github.com/Tiefseetauchner/EmoteSuggestor/pkg/LocalFunctions"
	"github.com/bwmarrin/discordgo"
)

var commands = []string{"a", "b"}
var listeningChannels = []string{}
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
		sendMessage(session, message.ChannelID, helpMsg)
		break
	case "suggest":
		if localFunctions.contains(listeningChannels, message.ChannelID)
		usage := "```\n" +
			"!suggest name [link] ...: suggests adding emote from link with name name to emojis\n" +
			"!suggest name: suggests adding emote from attachment with name name to emojis\n" +
			"```"
		if len(message.Attachments) > 0 {
			args := strings.Split(strings.TrimSpace(arguments), " ")
			if len(args) != 1 {
				sendMessage(session, message.ChannelID, "Invalid command. Usage: "+usage)
			}
			sendMessage(session, message.ChannelID, "Suggesting adding: "+args[0])
		} else if len(arguments) > 0 {
			args := strings.Split(strings.TrimSpace(arguments), " ")
			if len(args) != 2 {
				sendMessage(session, message.ChannelID, "Invalid command. Usage: "+usage)
			} else {
				sendMessage(session, message.ChannelID, "Suggesting adding: "+args[0])
			}
		} else {
			sendMessage(session, message.ChannelID, "No link or Embed provided. Usage:\n"+usage)
		}
	default:
		sendMessage(session, message.ChannelID, "Command "+command+" not found. Type !help for a list of commands")
	}
}

func suggestEmote(name string, link string) {

}

func sendMessage(session *discordgo.Session, id string, s string) {
	_, _ = session.ChannelMessageSend(id, s)
}
