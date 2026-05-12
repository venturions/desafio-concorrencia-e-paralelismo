package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/processor"
	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/utils"
)

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
	r := processor.ProcessSequential(files)
	elapsed := time.Since(start)

	r.PrintReport("Sequential", elapsed)

	// start2 := time.Now()
	// r2 := processor.ProcessConcurrentNaive(files)
	// elapsed2 := time.Since(start2)
	// r2.PrintReport("Concurrent naive", elapsed2)

	start3 := time.Now()
	r3 := processor.ProcessConcurrentMutex(files)
	elapsed3 := time.Since(start3)
	r3.PrintReport("Concurrent mutex", elapsed3)

	start4 := time.Now()
	r4 := processor.ProcessPipeline(files, runtime.NumCPU())
	elapsed4 := time.Since(start4)
	r4.PrintReport("Pipeline", elapsed4)

}
