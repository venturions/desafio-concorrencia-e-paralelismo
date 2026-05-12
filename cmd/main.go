package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/utils"
)

type Event struct {
	EventType string `json:"event_type"`
	Region    string `json:"region"`
}

type Report struct {
	totalEvents    int
	totalErrors    int
	eventsByType   map[string]int
	eventsByRegion map[string]int
}

func NewReport() *Report {
	return &Report{
		eventsByType:   make(map[string]int),
		eventsByRegion: make(map[string]int),
	}
}

func (r *Report) AddEvent(event Event) {
	r.totalEvents++
	r.eventsByType[event.EventType]++
	r.eventsByRegion[event.Region]++
}

func (r *Report) AddError() {
	r.totalErrors++
}

const dirPath = "../tmp/logs"

func main() {
	fmt.Println("Iniciando geração de arquivos de log...")

	if err := utils.GenerateMockFiles(dirPath, 100, 10000); err != nil {
		fmt.Printf("Erro ao gerar arquivos de log: %v\n", err)
		return
	}

	files, err := filepath.Glob("../tmp/logs/*.json")

	if err != nil {
		fmt.Printf("Erro ao buscar os arquivos gerados. %v\n", err)
		return
	}

	start := time.Now()
	r := ProcessSequential(files)
	elapsed := time.Since(start)

	fmt.Printf(
		"Total events: %d\n"+
			"Total errors: %d\n"+
			"Events by type: %v\n"+
			"Events by region: %v\n"+
			"Time elapsed: %v\n",
		r.totalEvents,
		r.totalErrors,
		r.eventsByType,
		r.eventsByRegion,
		elapsed,
	)

}

func ProcessSequential(files []string) *Report {
	r := NewReport()
	for _, file := range files {
		lines, err := utils.ReadFile(file)
		if err != nil {
			r.AddError()
			continue
		}

		for _, line := range lines {
			var event Event
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				r.AddError()
				continue
			}

			r.AddEvent(event)
		}
	}

	return r

}

func ProcessConcurrentNaive(files []string) *Report {
	r := NewReport()

	// TO DO IMPLEMENTATION

	return r

}

func ProcessConcurrentMutex(files []string) *Report {
	r := NewReport()

	// TO DO IMPLEMENTATION

	return r
}

func ProcessPipeline(files []string) *Report {
	r := NewReport()

	// TO DO IMPLEMENTATION

	return r
}
