package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"strings"
)

var prefix, channel, nickname string

func Init() {
	prefix = os.Getenv("PREFIX")
	channel = os.Getenv("CHANNEL")
	nickname = os.Getenv("NAME")
	LoadCommands()
}

func Start() {
	oauth := os.Getenv("SECRET")

	client := twitch.NewClient(nickname, oauth)

	client.OnConnect(func() {
		log.Println("Connected to " + channel)
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		onMessage(client, message)
	})

	client.Join(channel)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func onMessage(client *twitch.Client, message twitch.PrivateMessage) {
	if prefix == "" {
		log.Fatalf("No prefix defined")
	}

	if strings.HasPrefix(message.Message, prefix) {
		commandString := getCommandStringFromMessage(message)
		for _, command := range CommandList {
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
