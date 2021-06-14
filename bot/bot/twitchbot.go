package bot

import (
	"errors"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"strings"
)

var prefix, channel, nickname string
var messageCount = 0

func Init() {
	log.Println("Setting up bot...")
	prefix = os.Getenv("PREFIX")
	channel = os.Getenv("CHANNEL")
	nickname = os.Getenv("NAME")

	if prefix == "" {
		panic(errors.New("no PREFIX defined"))
	}
	if channel == "" {
		panic(errors.New("no CHANNEL defined"))
	}
	if nickname == "" {
		panic(errors.New("no NAME defined"))
	}

	LoadCommands()
}

func Start() {
	log.Println("Starting bot...")
	oauth := os.Getenv("SECRET")

	if oauth == "" {
		panic(errors.New("no SECRET given in .env"))
	}

	client := twitch.NewClient(nickname, oauth)

	client.OnConnect(func() {
		log.Println("Connected to " + channel)
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		onMessage(client, message)
	})

	client.Join(channel)

	log.Printf("Connecting to #%s...\n", channel)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func onMessage(client *twitch.Client, message twitch.PrivateMessage) {
	incrementMessageCount(message)

	handleIntervalMessage(client)

	if prefix == "" {
		log.Fatalf("No prefix defined")
	}

	if strings.HasPrefix(message.Message, prefix) {
		commandString := getCommandStringFromMessage(message)
		for _, command := range InvokableCommandList {
			if command.Invocation == commandString {
				if !command.ModOnly || (command.ModOnly && (message.User.Badges["moderator"] == 1) || (message.User.Badges["broadcaster"] == 1)) {
					client.Say(channel, command.Message)
				}
			}
		}
	}
}

func parseMessageText(message twitch.PrivateMessage) string {
	messageText := message.Message[len(prefix):]
	return strings.ToLower(messageText)
}

func getCommandStringFromMessage(message twitch.PrivateMessage) string {
	messageText := parseMessageText(message)
	return strings.Split(messageText, " ")[0]
}

func incrementMessageCount(message twitch.PrivateMessage) {
	if message.User.Name != nickname {
		messageCount += 1
	}
}

func handleIntervalMessage(client *twitch.Client) {
	for _, intervalMessage := range IntervalMessages {
		if messageCount % intervalMessage.MessageInterval == 0 {
			client.Say(channel, intervalMessage.Message)
		}
	}
}