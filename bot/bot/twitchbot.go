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
		err, command := parseCommand(message)
		if err != nil {
			if err.Error() == "no command detected" {
				log.Println("Prefix used in message but no valid command detected: " + message.Message)
			} else {
				log.Println(err.Error())
			}
		} else {
			log.Println(">>>>>>>>>>>>>>>>>>>>>")
			switch command {
			case PING:
				handlePing(client)
				break
			case DISCORD:
				handleDiscord(client)
				break
			case TWITTER:
				handleTwitter(client)
				break
			case INSTAGRAM:
				handleInstagram(client)
				break
			default:
				log.Println("Valid command detected but no implementation given")
			}
			log.Println("<<<<<<<<<<<<<<<<<<<<<")
		}
	}
}

func parseMessageText(message twitch.PrivateMessage) string {
	messageText := message.Message[len(prefix):]
	return strings.ToLower(messageText)
}

func parseCommand(message twitch.PrivateMessage) (error, Command) {
	messageText := parseMessageText(message)
	if strings.HasPrefix(messageText, string(PING)) {
		return nil, PING
	}
	if strings.HasPrefix(messageText, string(DISCORD)) {
		return nil, DISCORD
	}
	if strings.HasPrefix(messageText, string(TWITTER)) {
		return nil, TWITTER
	}
	if strings.HasPrefix(messageText, string(INSTAGRAM)) {
		return nil, INSTAGRAM
	}
	return errors.New("no command detected"), ""
}

func handlePing(client *twitch.Client) {
	log.Println("Ping command received")
	client.Say(channel, "Pong!")
	log.Println("Ping command handled")
}

func handleDiscord(client *twitch.Client) {
	log.Println("Discord command received")
	discordLink := os.Getenv("DISCORD")
	if discordLink != "" {
		client.Say(channel, discordLink)
	}
	log.Println("Discord command handled")
}

func handleTwitter(client *twitch.Client) {
	log.Println("Twitter command received")
	twitterLink := os.Getenv("TWITTER")
	if twitterLink != "" {
		client.Say(channel, twitterLink)
	}
	log.Println("Twitter command handled")
}

func handleInstagram(client *twitch.Client) {
	log.Println("Instagram command received")
	instagramLink := os.Getenv("INSTAGRAM")
	if instagramLink != "" {
		client.Say(channel, instagramLink)
	}
	log.Println("Instagram command handled")
}
