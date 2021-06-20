package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"math"
	"strconv"
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

func TestGetCommandStringFromMessage_SingleWordMessage(t *testing.T) {
	prefix = "!"
	err, result := getCommandStringFromMessage(twitch.PrivateMessage{Message: "!test"})
	if err != nil {
		t.Error("Unexpected error: " + err.Error())
	}
	if result != "test" {
		t.Error("Test Failed: Expected to be 'test' but was " + result)
	}
}

func TestGetCommandStringFromMessage_MultiWordMessage(t *testing.T) {
	prefix = "!"
	err, result := getCommandStringFromMessage(twitch.PrivateMessage{Message: "!test message"})
	if err != nil {
		t.Error("Unexpected error: " + err.Error())
	}
	if result != "test" {
		t.Error("Test Failed: Expected to be 'test' but was " + result)
	}
}

func TestGetCommandStringFromMessage_JustPrefix(t *testing.T) {
	prefix = "!"
	err, result := getCommandStringFromMessage(twitch.PrivateMessage{Message: "!"})
	if err == nil {
		t.Error("Expected error but received nil")
	}
	if result != "" {
		t.Error("Test Failed: Expected to be '' but was " + result)
	}
}

func TestIncrementMessageCount_NormalIncrement(t *testing.T) {
	messageCount = 0
	nickname = "test"
	incrementMessageCount(twitch.PrivateMessage{User: twitch.User{Name: "different"}})
	if messageCount != 1 {
		t.Error("Test Failed: Expected messageCount to be 1 but was " + strconv.FormatInt(int64(messageCount), 10))
	}
}

func TestIncrementMessageCount_MessageFromBot(t *testing.T) {
	messageCount = 0
	nickname = "test"
	incrementMessageCount(twitch.PrivateMessage{User: twitch.User{Name: "test"}})
	if messageCount != 0 {
		t.Error("Test Failed: Expected messageCount to be 0 but was " + strconv.FormatInt(int64(messageCount), 10))
	}
}

func TestIncrementMessageCount_MaxMessageCount(t *testing.T) {
	messageCount = math.MaxUint32 - 1
	nickname = "test"
	incrementMessageCount(twitch.PrivateMessage{User: twitch.User{Name: "different"}})
	if messageCount != 1 {
		t.Error("Test Failed: Expected messageCount to be 1 but was " + strconv.FormatInt(int64(messageCount), 10))
	}
}

func TestGetParametersFromMessage_MismatchingNumberOfParameters(t *testing.T) {
	prefix = "!"
	command := InvokableCommand{Parameters: []CommandParameter{{Name: "test"}, {Name: "another"}}}

	err, result := getParametersFromMessage(twitch.PrivateMessage{Message: "!command first"}, command)
	if err == nil {
		t.Error("Test Failed: Expected error but was nil")
	} else {
		if err.Error() != "number of parameters given does not match the number of parameters in the command" {
			t.Error("Test Failed: Expected error to be 'number of parameters given does not match the number of parameters in the command' but was: " + err.Error())
		}
	}

	if result != nil {
		t.Logf("Test Failed: Expected result to be nil but was: %s\n", result)
		t.Fail()
	}
}

func TestGetParametersFromMessage_ValidNumberOfParameters(t *testing.T) {
	prefix = "!"
	command := InvokableCommand{Parameters: []CommandParameter{{Name: "test"}, {Name: "another"}}}

	err, result := getParametersFromMessage(twitch.PrivateMessage{Message: "!command first second"}, command)
	if err != nil {
		t.Error("Test Failed: Expected no error but was: " + err.Error())
	}
	if !(len(result) == 2 || result[0] == "first" || result[1] == "second") {
		t.Logf("Test Failed: Incorrect slice returned, expected [\"first\", \"second\"] but received %s", result)
		t.Fail()
	}
}

func TestReplaceReservedKeywordsWithValues_ReplaceUsernameOnce(t *testing.T) {
	result := replaceReservedKeywordsWithValues("hello $username", twitch.PrivateMessage{User: twitch.User{Name: "testUsername"}})
	if result != "hello testUsername" {
		t.Error("Test Failed: Expected result to be 'hello testUsername' but was : " + result)
	}
}

func TestReplaceReservedKeywordsWithValues_ReplaceUsernameMultipleTimes(t *testing.T) {
	result := replaceReservedKeywordsWithValues("hello $username and $username", twitch.PrivateMessage{User: twitch.User{Name: "testUsername"}})
	if result != "hello testUsername and testUsername" {
		t.Error("Test Failed: Expected result to be 'hello testUsername and testUsername' but was : " + result)
	}
}

func TestReplaceCommandPlaceholdersWithValues_ValidOneReplacement(t *testing.T) {
	result := replaceCommandPlaceholdersWithValues("test $first", []CommandParameter{{Name: "first"}}, []string{"testValueOne"})
	if result != "test testValueOne" {
		t.Error("Test Failed: Expected result to be 'test testValueOne' but was: " + result)
	}
}

func TestReplaceCommandPlaceholdersWithValues_ValidMultipleReplacementsOfOneParameter(t *testing.T) {
	result := replaceCommandPlaceholdersWithValues("test $first and $first", []CommandParameter{{Name: "first"}}, []string{"testValueOne"})
	if result != "test testValueOne and testValueOne" {
		t.Error("Test Failed: Expected result to be 'test testValueOne and testValueOne' but was: " + result)
	}
}

func TestReplaceCommandPlaceholdersWithValues_ValidMultipleReplacementsOfMultipleParameters(t *testing.T) {
	result := replaceCommandPlaceholdersWithValues("test $first and $second and $first and $third", []CommandParameter{{Name: "first"}, {Name: "second"}, {Name: "third"}}, []string{"testValueOne", "testValueTwo", "testValueThree"})
	if result != "test testValueOne and testValueTwo and testValueOne and testValueThree" {
		t.Error("Test Failed: Expected result to be 'test testValueOne and testValueTwo and testValueOne and testValueThree' but was: " + result)
	}
}
