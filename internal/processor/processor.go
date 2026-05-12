package processor

import (
	"encoding/json"
	"sync"

	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/report"
	"github.com/venturions/desafio-concorrencia-e-paralelismo/internal/utils"
)

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

func ProcessPipeline(files []string) *report.Report {
	r := report.NewReport()

	// TO DO IMPLEMENTATION

	return r
}
