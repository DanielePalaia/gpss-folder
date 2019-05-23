package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {

	prop, _ := ReadPropertiesFile("./properties.ini")
	folder := prop["folder"]
	kafkaIp := prop["kafkaIp"]
	topicJson := prop["topicJson"]
	topicCsv := prop["topicCsv"]
	topicText := prop["topicTxt"]

	/* Remove last / if present*/
	if last := len(folder) - 1; last >= 0 && folder[last] == '/' {
		folder = folder[:last]
	}
	kafkaClient := makeKafkaEngine(kafkaIp, topicJson, topicCsv, topicText)

	/* Run the directory management object */
	folderEngine := makeFolderEngine(folder, kafkaClient)
	folderEngine.createOperationalFolders()
	folderEngine.folderListen()

}

type AppConfigProperties map[string]string

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
