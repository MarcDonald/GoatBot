package bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CommandParameter struct {
	Name string `json:"name"`
}

type InvokableCommand struct {
	Invocation string             `json:"invocation"`
	Parameters []CommandParameter `json:"parameters"`
	Message    string             `json:"message"`
	ModOnly    bool               `json:"mod_only"`
}

type IntervalMessage struct {
	Message         string `json:"message"`
	MessageInterval int    `json:"message_interval"`
}

var InvokableCommandList []InvokableCommand
var IntervalMessageList []IntervalMessage

// LoadCommands loads the commands from the commands/ folder into the bot's memory
func LoadCommands() {
	log.Println("Loading commands")

	var files []string

	filepathError := filepath.Walk("commands/", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if filepathError != nil {
		log.Fatalf(filepathError.Error())
	}

	for _, filePath := range files {
		err := loadFile(filePath)
		if err != nil {
			log.Println("Error loading file " + filePath + ": " + err.Error())
		}
	}

	log.Printf("%d invokable commands successfully loaded\n", len(InvokableCommandList))
	log.Printf("%d interval commands successfully loaded\n", len(IntervalMessageList))
}

// Loads an individual command file and store the command into memory
func loadFile(filePath string) error {
	currentFile, err := os.Open(filePath)
	if currentFile != nil {
		defer func(currentFile *os.File) {
			err := currentFile.Close()
			if err != nil {
				log.Println("Error closing " + filePath + ": " + err.Error())
			}
		}(currentFile)
	}
	if err != nil {
		return err
	}

	fileInfo, err := currentFile.Stat()
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		fileData, fileReadErr := ioutil.ReadFile(filePath)

		if fileReadErr != nil {
			log.Println("Error loading fileData " + filePath)
		}

		if strings.HasSuffix(filePath, ".interval.json") {
			commandFromFile := IntervalMessage{}
			_ = json.Unmarshal(fileData, &commandFromFile)
			IntervalMessageList = append(IntervalMessageList, commandFromFile)
		} else if strings.HasSuffix(filePath, ".command.json") {
			commandFromFile := InvokableCommand{}
			_ = json.Unmarshal(fileData, &commandFromFile)
			InvokableCommandList = append(InvokableCommandList, commandFromFile)
		} else {
			return errors.New("file does not have a valid suffix (i.e. `.command.json` or `.interval.json`")
		}
	}
	return nil
}
