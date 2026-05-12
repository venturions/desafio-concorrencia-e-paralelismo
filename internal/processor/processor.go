package processor

import (
	"encoding/json"
	"sync"

	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/report"
	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/utils"
)

type ProcessResult struct {
	Event report.Event
	Err   error
}

func ProcessSequential(files []string) *report.Report {
	r := report.NewReport()
	for _, file := range files {
		lines, err := utils.ReadFile(file)
		if err != nil {
			r.AddError()
			continue
		}

		for _, line := range lines {
			var event report.Event
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				r.AddError()
				continue
			}

			r.AddEvent(event)
		}
	}

	return r

}

func ProcessConcurrentNaive(files []string) *report.Report {
	r := report.NewReport()
	var wg = sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			lines, err := utils.ReadFile(file)

			if err != nil {
				r.AddError()
				return
			}

			for _, line := range lines {
				var event report.Event
				if err := json.Unmarshal([]byte(line), &event); err != nil {
					r.AddError()
					continue
				}

				r.AddEvent(event)
			}
		}(file)
	}

	wg.Wait()
	return r

}

func ProcessConcurrentMutex(files []string) *report.Report {
	r := report.NewReport()

	var wg = sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			lines, err := utils.ReadFile(file)

			if err != nil {
				r.AddErrorSafe()
				return
			}

			for _, line := range lines {
				var event report.Event
				if err := json.Unmarshal([]byte(line), &event); err != nil {
					r.AddErrorSafe()
					continue
				}

				r.AddEventSafe(event)
			}
		}(file)
	}

	wg.Wait()
	return r

}

func ProcessPipeline(files []string, numWorkers int) *report.Report {
	r := report.NewReport()
	if numWorkers < 1 {
		numWorkers = 1
	}

	jobs := make(chan string, len(files))
	results := make(chan ProcessResult, 1000)

	var workersWG sync.WaitGroup
	var aggregatorWG sync.WaitGroup

	aggregatorWG.Add(1)
	go func() {
		defer aggregatorWG.Done()

		for result := range results {
			if result.Err != nil {
				r.AddError()
				continue
			}

			r.AddEvent(result.Event)
		}
	}()

	worker := func() {
		defer workersWG.Done()

		for file := range jobs {
			lines, err := utils.ReadFile(file)
			if err != nil {
				results <- ProcessResult{Err: err}
				continue
			}

			for _, line := range lines {
				var event report.Event
				if err := json.Unmarshal([]byte(line), &event); err != nil {
					results <- ProcessResult{Err: err}
					continue
				}

				results <- ProcessResult{Event: event}
			}
		}
	}

	for i := 0; i < numWorkers; i++ {
		workersWG.Add(1)
		go worker()
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	workersWG.Wait()
	close(results)
	aggregatorWG.Wait()

	return r
}
