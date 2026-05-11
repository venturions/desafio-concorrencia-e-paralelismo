package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/utils"
)

type Report struct {
	totalEvents    int
	totalErrors    int
	eventsByType   map[string]int
	eventsByRegion map[string]int
}

func NewReport() *Report {
	return &Report{
		totalEvents: 0, totalErrors: 0, eventsByType: make(map[string]int), eventsByRegion: make(map[string]int),
	}
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

	fmt.Printf("Arquivos de logs gerados com sucesso.\n %s", files)
}

func ProcessSequential(files []string) *Report {
	start := time.Now()
	elapsed := time.Since(start)

	r := NewReport()

	// TO DO IMPLEMENTATION

	return r

}

func ProcessConcurrentNaive(files []string) *Report {
	start := time.Now()
	elapsed := time.Since(start)
	r := NewReport()

	// TO DO IMPLEMENTATION

	return r

}

func ProcessConcurrentMutex(files []string) *Report {
	start := time.Now()
	elapsed := time.Since(start)
	r := NewReport()

	// TO DO IMPLEMENTATION

	return r
}

func ProcessPipeline(files []string) *Report {
	start := time.Now()
	elapsed := time.Since(start)
	r := NewReport()

	// TO DO IMPLEMENTATION

	return r
}
