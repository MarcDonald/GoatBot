package bot

import (
	"errors"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"strings"
)

var prefix, channel string

func Init() {
	prefix = os.Getenv("PREFIX")
}

func Start() {
	nickname := os.Getenv("NAME")
	channel = os.Getenv("CHANNEL")
	oauth := os.Getenv("PASSWORD")

	client := twitch.NewClient(nickname, oauth)

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

	if strings.HasPrefix(message.Message, prefix+" ") {
		err, command := parseCommand(message)
		if err != nil {
			if err.Error() == "no command detected" {
				log.Println("Prefix used in message but no valid command detected: " + message.Message)
			} else {
				log.Println(err.Error())
			}
		} else {
			switch command {
			case PING:
				handlePing(client)
				break
			default:
				log.Println("Valid command detected but no implementation given")
			}
		}
	}
}

func parseMessageText(message twitch.PrivateMessage) string {
	messageText := message.Message[len(prefix)+1:]
	return strings.ToLower(messageText)
}

func parseCommand(message twitch.PrivateMessage) (error, Command) {
	messageText := parseMessageText(message)
	if strings.HasPrefix(messageText, string(PING)) {
		return nil, PING
	}
	return errors.New("no command detected"), ""
}

func handlePing(client *twitch.Client) {
	log.Println("Ping command received")
	client.Say(channel, "Pong!")
	log.Println("Ping command handled")
}