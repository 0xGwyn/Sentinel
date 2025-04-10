package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler struct {
	scheduler gocron.Scheduler
	// coordinator *Coordinator
	config Config
}

func NewScheduler(config Config) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %v", err)
	}

	return &Scheduler{
		scheduler: scheduler,
		// coordinator: NewCoordinator(),
		config: config,
	}, nil
}

func (s *Scheduler) Start() error {
	for _, jobType := range []JobType{SubfinderJob, DnsxJob, HttpxJob} {
		var jobDuration int
		var taskLogic any

		if jobType == SubfinderJob {
			jobDuration = s.config.SubfinderInterval
			taskLogic = subfinderTask
		} else if jobType == DnsxJob {
			jobDuration = s.config.DnsxInterval
			taskLogic = dnsxTask
		} else if jobType == HttpxJob {
			jobDuration = s.config.HttpxInterval
			taskLogic = httpxTask
		}

		_, err := s.scheduler.NewJob(
			gocron.DurationJob(
				time.Duration(jobDuration)*time.Second,
			),
			gocron.NewTask(
				taskLogic,
			),
			gocron.WithName(string(jobType)+"-job"),
			gocron.WithEventListeners(
				gocron.BeforeJobRunsSkipIfBeforeFuncErrors(func(jobID uuid.UUID, jobName string) error {
					// if s.coordinator.CanRun(jobType) {
					log.Printf("before run: %s", jobName)
					// err := s.coordinator.StartJob(jobType)
					// if err != nil {
					// 	return err // Skip job if StartJob fails
					// }
					return nil // Allow job to run
					// }
					// return fmt.Errorf("cannot run job(%s) - previous job still running", jobType) // Skip job
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					// s.coordinator.EndJob(jobType)
					log.Printf("after run Completed %s", jobName)
				}),
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					// s.coordinator.EndJob(jobType)
					log.Printf("Inside after run Completed with error %s", jobName)
				}),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to create job %s: %v", jobType, err)
		}

	}

	s.scheduler.Start()

	return nil
}

func (s *Scheduler) Stop() {
	if err := s.scheduler.Shutdown(); err != nil {
		log.Printf("Error shutting down scheduler: %v", err)
	}

}

func subfinderTask() error {
	log.Println("inside subfinder")
	// return fmt.Errorf("error testing")
	return nil
	// Subfinder job
}

func httpxTask() error {
	log.Println("inside httpx")
	// Httpx job
	return nil
}

func dnsxTask() error {
	log.Println("inside dnsx")
	// Dnsx job
	return nil
}
