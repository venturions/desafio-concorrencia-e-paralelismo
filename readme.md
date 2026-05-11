# Desafio: Concorrência e Paralelismo em Go

Este desafio simula um sistema de auditoria de logs que leva 45 minutos para processar dados. O objetivo é criar um motor de agregação em Go demonstrando domínio progressivo de concorrência por meio de quatro implementações:

- sequencial
- concorrente ingênua, com race conditions
- concorrente com `Mutex`
- worker pool com channels

O foco é aplicar conceitos essenciais para sistemas de alta performance, comparando segurança, concorrência e escalabilidade na prática.

## Contexto

Em um cenário hipotético, uma equipe de Business Intelligence precisa de um relatório de auditoria com todos os eventos da última hora. O problema é que o job atual de processamento leva 45 minutos para rodar.

Sua missão é criar do zero um novo motor de agregação em Go para eliminar esse gargalo. Os servidores despejam milhares de arquivos de log JSON em um diretório, e você deve usar goroutines, channels e mutexes para implementar esse job de agregação.

## Objetivos

O objetivo é demonstrar domínio dos seguintes conceitos:

- Processamento síncrono, concorrência e paralelismo
- Detecção e compreensão de `race conditions`
- Uso correto de `sync.Mutex` e análise de contenção de lock
- Pipelines escaláveis com `Fan-Out` e `Fan-In`
- Worker pools para controle de paralelismo
- Ciclo de vida robusto com `sync.WaitGroup` e fechamento de channels

## Estrutura, regras e requisitos

### Requisitos gerais

- Criar um arquivo `main.go`
- Implementar as 4 funções descritas abaixo
- Na função `main`, você deve:
  - gerar um conjunto de arquivos de log de teste usando `GenerateMockFiles`
  - chamar cada uma das 4 funções
  - medir o tempo de execução com `time.Since(start)`
  - imprimir os resultados para comparação

## Tarefas

### Parte 1: Baseline

Implementar a versão mais simples e direta, processando um arquivo por vez, de forma sequencial, na goroutine principal.

**Assinatura**

```go
func ProcessSequential(files []string) *Report
```

**Lógica**

- Criar um `Report` com `NewReport()`
- Iterar com `for...range` sobre `files`
- Abrir cada arquivo e ler linha por linha com `bufio.Scanner`
- Para cada linha válida, decodificar o JSON para um `Event` e chamar `report.AddEvent(event)`
- Para cada linha inválida, chamar `report.AddError()`

### Parte 2: Concorrência ingênua

Processar cada arquivo em sua própria goroutine. Esta implementação é intencionalmente falha.

**Assinatura**

```go
func ProcessConcurrentNaive(files []string) *Report
```

**Lógica**

- Criar um único `Report` compartilhado
- Criar um `sync.WaitGroup`
- Para cada arquivo:
  - chamar `wg.Add(1)`
  - iniciar uma goroutine `go func(filename string) { ... }`
  - dentro da goroutine, usar `defer wg.Done()`
  - abrir, ler e decodificar o arquivo
  - para cada evento, chamar `report.AddEvent(event)`
  - para cada erro, chamar `report.AddError()`
- Aguardar com `wg.Wait()`
- Retornar o relatório

### Parte 3: Corrigindo com `Mutex`

Corrigir a versão concorrente usando memória compartilhada protegida.

**Assinatura**

```go
func ProcessConcurrentMutex(files []string) *Report
```

**Lógica**

- Copiar a implementação da Parte 2
- Manter `WaitGroup` e goroutines
- Dentro da goroutine, usar as versões seguras:
  - `report.AddEventSafe(event)`
  - `report.AddErrorSafe()`

### Parte 4: Worker Pool e Channels

Refatorar para a arquitetura mais escalável, evitando memória compartilhada e seguindo o lema do Go:

> Não comunique compartilhando memória; compartilhe memória comunicando.

**Assinatura**

```go
func ProcessPipeline(files []string, numWorkers int) *Report
```

**Lógica**

- Criar dois channels:
  - `jobs := make(chan string, len(files))`
  - `results := make(chan ProcessResult, 1000)`
- Usar uma struct `ProcessResult` que possa conter um `Event` ou um `error`
- **Workers (`Fan-Out`)**
  - criar um `sync.WaitGroup` para os workers
  - iniciar `numWorkers` goroutines
  - cada worker lê do canal `jobs`
  - o worker processa o arquivo e envia `ProcessResult` para `results`
- **Aggregator (`Fan-In`)**
  - criar um `Report` com `NewReport()`
  - iniciar uma goroutine agregadora que será a única leitora de `results`
  - para cada resultado:
    - se houver erro, chamar `report.AddError()`
    - caso contrário, chamar `report.AddEvent(res.Event)`
- **Coordenação**
  - enviar todos os arquivos para `jobs`
  - chamar `close(jobs)`
  - aguardar os workers com `wg.Wait()`
  - chamar `close(results)`
  - aguardar o agregador terminar
  - retornar o `report`

## Resultados esperados

### Correção dos dados

- As implementações das Partes 1, 3 e 4 devem produzir um relatório final idêntico
- A Parte 2 deve falhar em produzir o resultado correto ou, no mínimo, ser comprovadamente insegura

### Robustez e segurança

- `go run -race .` deve passar sem avisos nas Partes 1, 3 e 4
- `go run -race .` deve detectar `DATA RACE` na Parte 2

### Performance mensurável

- Medir desempenho com `time.Since(start)` ou com o pacote `testing`
- Comprovar melhoria das Partes 3 e 4 em relação à baseline da Parte 1
- Usar um volume significativo de dados, por exemplo:
  - 100 arquivos
  - 10.000 linhas por arquivo

### Adesão aos padrões exigidos

- A Parte 3 deve usar `sync.Mutex` ou `sync.RWMutex`
- A Parte 4 deve usar obrigatoriamente worker pool e pipeline de channels com `Fan-Out/Fan-In`
- O agregador da Parte 4 não deve precisar de lock

## Checklist de entrega

- [ ] Integrar a função `GenerateMockFiles()`
- [ ] Implementar `ProcessSequential()` e medir tempo
- [ ] Implementar `ProcessConcurrentNaive()` e confirmar `DATA RACE` com `go run -race .`
- [ ] Implementar `ProcessConcurrentMutex()` e validar ausência de race conditions
- [ ] Implementar `ProcessPipeline()` com worker pool e channels
- [ ] Criar a função `main()` chamando as 4 implementações
- [ ] Gerar dados de teste e comparar performance
- [ ] Validar que as Partes 1, 3 e 4 produzem resultados idênticos
