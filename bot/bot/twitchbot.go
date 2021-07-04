package bot

import (
	"errors"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"strings"
)

var prefix, channel, nickname string

type ChatClient interface {
	Say(channel, text string)
}

// Init initializes variables for the bot and loads the commands
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

// Start starts the bot
func Start() {
	log.Println("Starting bot...")
	oauth := os.Getenv("SECRET")

	if oauth == "" {
		panic(errors.New("no SECRET given in .env"))
	}

	client := twitch.NewClient(nickname, oauth)
	commandHandler := CommandHandler{}

	client.OnConnect(func() {
		log.Println("Connected to " + channel)
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		onMessage(&commandHandler, client, message)
	})

	client.Join(channel)

	log.Printf("Connecting to #%s...\n", channel)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

// TODO test
// Handle message event
func onMessage(handler CommandProcessor, client ChatClient, message twitch.PrivateMessage) {
	handler.IncrementMessageCount(message)
	handler.HandleIntervalMessage(client)

	if prefix == "" {
		log.Fatalf("No prefix defined")
	}

	if strings.HasPrefix(message.Message, prefix) {
		err, commandString := handler.GetCommandStringFromMessage(message)
		if err != nil {
			log.Println("Error parsing command from message: " + err.Error())
		} else {
			for _, command := range InvokableCommandList {
				if handler.HasCommandBeenInvoked(command, commandString) {
					if handler.HasPermissionToInvoke(command, message) {
						formattedMessage := handler.ReplaceReservedKeywordsWithValues(command.Message, message)
						if len(command.Parameters) != 0 {
							err, messageParameters := handler.GetParametersFromMessage(message, command)
							if err != nil {
								client.Say(channel, "Invalid usage of command")
								log.Println(err.Error())
							} else {
								formattedMessage = handler.ReplaceCommandPlaceholdersWithValues(formattedMessage, command.Parameters, messageParameters)
								client.Say(channel, formattedMessage)
							}
						} else {
							client.Say(channel, formattedMessage)
						}
					}
				}
			}
		}
	}
}
