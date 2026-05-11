package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Event struct {
	EventType string `json:"event_type"`
	Region    string `json:"region"`
}

// GenerateMockFiles cria arquivos de log JSON para serem processados. Essa função deve ser utilizada para testar o processamento de logs.
func GenerateMockFiles(dir string, numFiles, eventsPerFile int) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	eventTypes := []string{"click", "view", "purchase", "login"}
	regions := []string{"us-east-1", "eu-west-1", "ap-southeast-2", "sa-east-1"}

	for i := 0; i < numFiles; i++ {
		filePath := filepath.Join(dir, fmt.Sprintf("log_%03d.json", i))
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}

		for j := 0; j < eventsPerFile; j++ {
			var line string
			if j%50 == 0 && j > 0 {
				line = "this is not valid json\n"
			} else {
				event := Event{
					EventType: eventTypes[(i+j)%len(eventTypes)],
					Region:    regions[j%len(regions)],
				}
				data, err := json.Marshal(event)
				if err != nil {
					file.Close()
					return err
				}
				line = string(data) + "\n"
			}

			if _, err := file.WriteString(line); err != nil {
				file.Close()
				return fmt.Errorf("write %s: %w", filePath, err)
			}
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}
