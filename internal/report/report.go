package report

import (
	"fmt"
	"sync"
	"time"
)

type Report struct {
	mu             sync.Mutex
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

func (r *Report) AddEventSafe(event Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.AddEvent(event)

}

func (r *Report) AddError() {
	r.totalErrors++
}

func (r *Report) AddErrorSafe() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.AddError()

}

func (r *Report) PrintReport(label string, elapsed time.Duration) {
	fmt.Printf(
		"\n[%s]\n"+
			"Total events: %d\n"+
			"Total errors: %d\n"+
			"Events by type: %v\n"+
			"Events by region: %v\n"+
			"Time elapsed: %v\n",
		label,
		r.totalEvents,
		r.totalErrors,
		r.eventsByType,
		r.eventsByRegion,
		elapsed,
	)
}
