package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"testing"
)

func TestHasPermissionToInvoke_NotModOnlyCommand(t *testing.T) {
	testCommand := InvokableCommand{
		ModOnly: false,
	}

	testMessage := twitch.PrivateMessage{
		User: twitch.User{},
	}
	result := hasPermissionToInvoke(testCommand, testMessage)

	if !result {
		t.Error("Test Failed: Expected to have permission to invoke command when not a mod only command")
	}
}

func TestHasPermissionToInvoke_ModOnlyCommandNotMod(t *testing.T) {
	testCommand := InvokableCommand{
		ModOnly: true,
	}

	testMessage := twitch.PrivateMessage{
		User: twitch.User{
			Badges: nil,
		},
	}
	result := hasPermissionToInvoke(testCommand, testMessage)

	if result {
		t.Error("Test Failed: Expected to not have permission to invoke command when mod only command invoked by non mod or broadcaster")
	}
}

func TestHasPermissionToInvoke_ModOnlyCommandBroadcaster(t *testing.T) {
	testCommand := InvokableCommand{
		ModOnly: true,
	}

	var testMessage = twitch.PrivateMessage{
		User: twitch.User{Badges: map[string]int{
			"broadcaster": 1,
		}},
	}
	result := hasPermissionToInvoke(testCommand, testMessage)

	if !result {
		t.Error("Test Failed: Expected to have permission to invoke command when not a mod only command")
	}
}

func TestHasPermissionToInvoke_ModOnlyCommandModerator(t *testing.T) {
	testCommand := InvokableCommand{
		ModOnly: true,
	}

	var testMessage = twitch.PrivateMessage{
		User: twitch.User{Badges: map[string]int{
			"moderator": 1,
		}},
	}
	result := hasPermissionToInvoke(testCommand, testMessage)

	if !result {
		t.Error("Test Failed: Expected to have permission to invoke command when not a mod only command")
	}
}

func TestParseMessageText_SingleCharacterPrefixMultiWordMessage(t *testing.T) {
	prefix = "!"
	testMessage := "!this is a test message"
	result := parseMessageText(twitch.PrivateMessage{Message: testMessage})
	if result != "this is a test message" {
		t.Error("Test Failed: Expected to be 'this is a test message' but was " + result)
	}
}

func TestParseMessageText_SingleCharacterPrefixSingleWordMessage(t *testing.T) {
	prefix = "!"
	testMessage := "!this"
	result := parseMessageText(twitch.PrivateMessage{Message: testMessage})
	if result != "this" {
		t.Error("Test Failed: Expected to be 'this' but was " + result)
	}
}

func TestParseMessageText_MultiCharacterPrefixMultiWordMessage(t *testing.T) {
	prefix = "prefix "
	testMessage := "prefix this is a test"
	result := parseMessageText(twitch.PrivateMessage{Message: testMessage})
	if result != "this is a test" {
		t.Error("Test Failed: Expected to be 'this is a test' but was " + result)
	}
}

func TestParseMessageText_SingleCharacterPrefixMixCases(t *testing.T) {
	prefix = "!"
	testMessage := "!this is a TEST"
	result := parseMessageText(twitch.PrivateMessage{Message: testMessage})
	if result != "this is a test" {
		t.Error("Test Failed: Expected to be 'this is a test' but was " + result)
	}
}

func TestParseMessageText_JustPrefix(t *testing.T) {
	prefix = "!"
	testMessage := "!"
	result := parseMessageText(twitch.PrivateMessage{Message: testMessage})
	if result != "" {
		t.Error("Test Failed: Expected to be '' but was " + result)
	}
}
