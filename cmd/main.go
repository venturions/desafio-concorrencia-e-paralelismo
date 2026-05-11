package main

import (
	"fmt"

	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/utils"
)

const dirPath = "../tmp/logs"

func main() {
	fmt.Println("Iniciando geração de arquivos de log...")

	if err := utils.GenerateMockFiles(dirPath, 10, 200); err != nil {
		fmt.Printf("Erro ao gerar arquivos de log: %v\n", err)
		return
	}

	fmt.Println("Arquivos de logs gerados com sucesso.")
}
