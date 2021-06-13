package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Command struct {
	Invocation string  `json:"invocation"`
	Message    string  `json:"message"`
	ModOnly    bool    `json:"mod_only"`
	Timer      float32 `json:"timer"`
}

var CommandList []Command

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

	log.Printf("%d commands successfully loaded\n", len(CommandList))
}

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

		commandFromFile := Command{}
		_ = json.Unmarshal(fileData, &commandFromFile)

		CommandList = append(CommandList, commandFromFile)
	}
	return nil
}
