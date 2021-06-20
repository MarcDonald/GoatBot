package bot

import (
	"errors"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"math"
	"os"
	"strings"
)

var prefix, channel, nickname string
var messageCount uint32 = 0

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

// TODO to test this probably need to do some DI
// Handle message event
func onMessage(client *twitch.Client, message twitch.PrivateMessage) {
	incrementMessageCount(message)

	handleIntervalMessage(client)

	if prefix == "" {
		log.Fatalf("No prefix defined")
	}

	if strings.HasPrefix(message.Message, prefix) {
		err, commandString := getCommandStringFromMessage(message)
		if err != nil {
			log.Println("Error parsing command from message: " + err.Error())
		} else {
			for _, command := range InvokableCommandList {
				if command.Invocation == commandString {
					if hasPermissionToInvoke(command, message) {
						if len(command.Parameters) != 0 {
							err, messageParameters := getParametersFromMessage(message, command)
							if err != nil {
								client.Say(channel, "Invalid usage of command")
								log.Println(err.Error())
							} else {
								formattedMessage := replaceReservedKeywordsWithValues(command.Message, message)
								formattedMessage = replaceCommandPlaceholdersWithValues(formattedMessage, command.Parameters, messageParameters)
								client.Say(channel, formattedMessage)
							}
						} else {
							client.Say(channel, command.Message)
						}
					}
				}
			}
		}
	}
}

// Returns true if a command is mod only and the user invoking the command is a mod or broadcaster, or if the command is not mod only
func hasPermissionToInvoke(command InvokableCommand, message twitch.PrivateMessage) bool {
	return !command.ModOnly || (command.ModOnly && (message.User.Badges["moderator"] == 1) || (message.User.Badges["broadcaster"] == 1))
}

// Returns the content of the message without the prefix
func parseMessageText(message twitch.PrivateMessage) string {
	messageText := message.Message[len(prefix):]
	return strings.ToLower(messageText)
}

// Returns the command string used to invoke the command
func getCommandStringFromMessage(message twitch.PrivateMessage) (error, string) {
	messageText := parseMessageText(message)
	command := strings.Split(messageText, " ")[0]
	if command != "" {
		return nil, command
	} else {
		return errors.New("missing command"), ""
	}
}

// Increments the message count (excluding messages from the bot)
func incrementMessageCount(message twitch.PrivateMessage) {
	if messageCount == math.MaxUint32-1 {
		messageCount = 0
	}

	if message.User.Name != nickname {
		messageCount += 1
	}
}

// Goes through the IntervalMessageList and sends a message if it is time to send that message
func handleIntervalMessage(client *twitch.Client) {
	for _, intervalMessage := range IntervalMessageList {
		if messageCount%uint32(intervalMessage.MessageInterval) == uint32(0) {
			// TODO to test this probably need to do some DI
			client.Say(channel, intervalMessage.Message)
		}
	}
}

// Returns the parameters used when invoking a command
func getParametersFromMessage(message twitch.PrivateMessage, command InvokableCommand) (error, []string) {
	var numParameters = len(command.Parameters)
	messageText := parseMessageText(message)
	messageWords := strings.Split(messageText, " ")

	if len(messageWords[1:]) < numParameters {
		return errors.New("number of parameters given does not match the number of parameters in the command"), nil
	}

	parameters := messageWords[1 : numParameters+1]
	return nil, parameters
}

// Returns the message with the reserved keywords replaced with their values
func replaceReservedKeywordsWithValues(commandMessage string, message twitch.PrivateMessage) string {
	var formattedMessage = commandMessage
	formattedMessage = strings.Replace(formattedMessage, "$username", message.User.Name, -1)
	return formattedMessage
}

// Returns the command message with the placeholders replaced with the given values
func replaceCommandPlaceholdersWithValues(commandMessage string, parameters []CommandParameter, messageParameters []string) string {
	var formattedMessage = commandMessage
	for i, parameter := range parameters {
		formattedMessage = strings.Replace(formattedMessage, "$"+parameter.Name, messageParameters[i], -1)
	}
	return formattedMessage
}
