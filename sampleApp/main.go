package main

import (
	log "sample_app/logFolder"
	"sample_app/metrics"
)

func main() {
	log.Info("Starting...")
	m, err := metrics.New("us-east-1")
	if err != nil {
		log.Errorf("Error initializing metrics: %s", err)
	}
	defer m.Publish()

	// incrementing counter 1 value till 2
	m.IncrementCounter1() // counter1 = 1
	m.IncrementCounter1() // counter1 = 2
}
