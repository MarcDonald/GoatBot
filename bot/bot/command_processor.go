package bot

import (
	"errors"
	"github.com/gempir/go-twitch-irc/v2"
	"math"
	"strings"
)

var messageCount uint32 = 0

type CommandProcessor interface {
	IncrementMessageCount(message twitch.PrivateMessage)
	HandleIntervalMessage(client ChatClient)
	HasPermissionToInvoke(command InvokableCommand, message twitch.PrivateMessage) bool
	GetCommandStringFromMessage(message twitch.PrivateMessage) (error, string)
	GetParametersFromMessage(message twitch.PrivateMessage, command InvokableCommand) (error, []string)
	ReplaceReservedKeywordsWithValues(commandMessage string, message twitch.PrivateMessage) string
	ReplaceCommandPlaceholdersWithValues(commandMessage string, parameters []CommandParameter, messageParameters []string) string
}

type CommandHandler struct{}

// IncrementMessageCount increments the message count (excluding messages from the bot)
func (h *CommandHandler) IncrementMessageCount(message twitch.PrivateMessage) {
	if messageCount == math.MaxUint32-1 {
		messageCount = 0
	}

	if message.User.Name != nickname {
		messageCount += 1
	}
}

// HandleIntervalMessage goes through the IntervalMessageList and sends a message if it is time to send that message
func (h *CommandHandler) HandleIntervalMessage(client ChatClient) {
	for _, intervalMessage := range IntervalMessageList {
		if messageCount%uint32(intervalMessage.MessageInterval) == uint32(0) {
			client.Say(channel, intervalMessage.Message)
		}
	}
}

// HasPermissionToInvoke returns true if a command is mod only and the user invoking the command is a mod or broadcaster, or if the command is not mod only
func (h *CommandHandler) HasPermissionToInvoke(command InvokableCommand, message twitch.PrivateMessage) bool {
	return !command.ModOnly || (command.ModOnly && (message.User.Badges["moderator"] == 1) || (message.User.Badges["broadcaster"] == 1))
}

// GetCommandStringFromMessage returns the command string used to invoke the command
func (h *CommandHandler) GetCommandStringFromMessage(message twitch.PrivateMessage) (error, string) {
	messageText := parseMessageText(message)
	command := strings.Split(messageText, " ")[0]
	if command != "" {
		return nil, command
	} else {
		return errors.New("missing command"), ""
	}
}

// GetParametersFromMessage returns the parameters used when invoking a command
func (h *CommandHandler) GetParametersFromMessage(message twitch.PrivateMessage, command InvokableCommand) (error, []string) {
	var numParameters = len(command.Parameters)
	messageText := parseMessageText(message)
	messageWords := strings.Split(messageText, " ")

	if len(messageWords[1:]) < numParameters {
		return errors.New("number of parameters given does not match the number of parameters in the command"), nil
	}

	parameters := messageWords[1 : numParameters+1]
	return nil, parameters
}

// ReplaceReservedKeywordsWithValues returns the message with the reserved keywords replaced with their values
func (h *CommandHandler) ReplaceReservedKeywordsWithValues(commandMessage string, message twitch.PrivateMessage) string {
	var formattedMessage = commandMessage
	formattedMessage = strings.Replace(formattedMessage, "$username", message.User.Name, -1)
	return formattedMessage
}

// ReplaceCommandPlaceholdersWithValues returns the command message with the placeholders replaced with the given values
func (h *CommandHandler) ReplaceCommandPlaceholdersWithValues(commandMessage string, parameters []CommandParameter, messageParameters []string) string {
	var formattedMessage = commandMessage
	for i, parameter := range parameters {
		formattedMessage = strings.Replace(formattedMessage, "$"+parameter.Name, messageParameters[i], -1)
	}
	return formattedMessage
}

// Returns the content of the message without the prefix
func parseMessageText(message twitch.PrivateMessage) string {
	messageText := message.Message[len(prefix):]
	return strings.ToLower(messageText)
}
