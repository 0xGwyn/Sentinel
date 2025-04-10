package scheduler

import (
	"log"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	log.Println("Test started")

	// Create test config with short intervals
	config := Config{
		SubfinderInterval: 5,
		HttpxInterval:     5,
		DnsxInterval:      100,
	}

	// Create new scheduler
	s, err := NewScheduler(config)
	if err != nil {
		t.Fatalf("Failed to create scheduler: %v", err)
	}

	// Start scheduler
	if err := s.Start(); err != nil {
		t.Fatalf("Failed to start scheduler: %v", err)
	}

	// Let it run for a while to observe multiple job executions
	time.Sleep(30 * time.Second)

	// Stop the scheduler
	s.Stop()
}
