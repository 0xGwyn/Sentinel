package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/0xgwyn/sentinel/database"
	"github.com/0xgwyn/sentinel/models"
	"github.com/0xgwyn/sentinel/modules"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler struct {
	scheduler   gocron.Scheduler
	coordinator *Coordinator
	config      Config
}

func NewScheduler(config Config) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %v", err)
	}

	return &Scheduler{
		scheduler:   scheduler,
		coordinator: NewCoordinator(),
		config:      config,
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
				time.Duration(jobDuration)*time.Minute,
			),
			gocron.NewTask(
				taskLogic,
			),
			gocron.WithName(string(jobType)+"-job"),
			gocron.WithEventListeners(
				gocron.BeforeJobRunsSkipIfBeforeFuncErrors(func(jobID uuid.UUID, jobName string) error {
					if s.coordinator.CanRun(jobType) {
						err := s.coordinator.StartJob(jobType)
						if err != nil {
							// Skip job if StartJob fails
							return err
						}
					}
					// Skip job
					return fmt.Errorf("cannot run job(%s) - previous job still running", jobType)
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					s.coordinator.EndJob(jobType, nil)
				}),
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					s.coordinator.EndJob(jobType, err)
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

func (s *Scheduler) Stop() error {
	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("error shutting down scheduler: %v", err)
	}

	return nil
}

func subfinderTask() error {
	log.Println("Running subfinder task")

	// Get domains collection
	domainsColl := database.GetDBCollection("domains")
	subdomainsColl := database.GetDBCollection("subdomains")

	// Find all domains
	cursor, err := domainsColl.Find(context.Background(), bson.M{})
	if err != nil {
		return fmt.Errorf("failed to fetch domains: %v", err)
	}
	defer cursor.Close(context.Background())

	// Iterate over domains
	for cursor.Next(context.Background()) {
		var domain models.Domain
		if err := cursor.Decode(&domain); err != nil {
			log.Printf("failed to decode domain: %v", err)
			continue
		}

		// Run subfinder for each domain
		results, err := modules.RunSubfinder(domain.Name)
		if err != nil {
			log.Printf("subfinder failed for domain %s: %v", domain.Name, err)
			continue
		}

		// Process each subdomain found by subfinder
		for _, result := range results {

			// Prepare filter for checking existing subdomain
			filter := bson.M{
				"domain": domain.Name,
				"name":   result.Subdomain,
			}

			// Try to find existing subdomain
			var existingSubdomain models.Subdomain
			err := subdomainsColl.FindOne(context.Background(), filter).Decode(&existingSubdomain)

			now := time.Now()

			if err == mongo.ErrNoDocuments {
				// Subdomain doesn't exist, create new one
				newSubdomain := models.Subdomain{
					Domain:    domain.Name,
					Name:      result.Subdomain,
					CreatedAt: bson.NewDateTimeFromTime(now),
					UpdatedAt: bson.NewDateTimeFromTime(now),
					Providers: result.Provider,
					WatchHTTP: true,
					WatchDNS:  true,
					DNSStatus: models.FreshSubdomain,
				}

				_, err = subdomainsColl.InsertOne(context.Background(), newSubdomain)
				if err != nil {
					log.Printf("failed to insert new subdomain %s: %v", result.Subdomain, err)
				}
			} else if err == nil {
				// Subdomain exists, check for new providers
				newProviders := make([]string, 0)
				existingProviders := make(map[string]bool)

				// Create map of existing providers
				for _, provider := range existingSubdomain.Providers {
					existingProviders[provider] = true
				}

				// Check for new providers
				for _, provider := range result.Provider {
					if !existingProviders[provider] {
						newProviders = append(newProviders, provider)
					}
				}

				// If new providers found, update the subdomain
				if len(newProviders) > 0 {
					update := bson.M{
						"$set": bson.M{
							"updated_at": bson.NewDateTimeFromTime(now),
						},
						"$push": bson.M{
							"providers": bson.M{
								"$each": newProviders,
							},
						},
					}

					_, err = subdomainsColl.UpdateOne(context.Background(), filter, update)
					if err != nil {
						log.Printf("failed to update subdomain %s providers: %v", result.Subdomain, err)
					}
				}
			} else {
				log.Printf("error checking subdomain %s: %v", result.Subdomain, err)
			}
		}
	}

	return nil
}

func dnsxTask() error {
	log.Println("Running dnsx task")

	// Get collections
	subdomainsColl := database.GetDBCollection("subdomains")
	dnsColl := database.GetDBCollection("dns")

	// Find all subdomains with WatchDNS true
	filter := bson.M{"watch_dns": true}
	cursor, err := subdomainsColl.Find(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to fetch subdomains: %v", err)
	}
	defer cursor.Close(context.Background())

	// Collect all subdomain names to process
	var subdomains []models.Subdomain
	if err = cursor.All(context.Background(), &subdomains); err != nil {
		return fmt.Errorf("failed to decode subdomains: %v", err)
	}

	if len(subdomains) == 0 {
		return nil
	}

	// Extract subdomain names for dnsx
	var subdomainNames []string
	for _, sub := range subdomains {
		subdomainNames = append(subdomainNames, sub.Name)
	}

	// Run dnsx
	dnsResults, err := modules.RunDnsx(subdomainNames, []string{"a", "aaaa", "cname", "ns", "ptr", "mx", "txt"}, 25)
	if err != nil {
		return fmt.Errorf("failed to run dnsx: %v", err)
	}

	now := time.Now()

	// Process results for each subdomain
	subdomainMap := make(map[string]models.Subdomain)
	for _, sub := range subdomains {
		subdomainMap[sub.Name] = sub
	}

	// Process results for each subdomain
	for _, result := range dnsResults {
		// Find the corresponding subdomain from our list
		currentSubdomain, exists := subdomainMap[result.Domain]
		if !exists {
			log.Printf("subdomain %s not found in our list", result.Domain)
			continue
		}

		// Check if there's any previous DNS record
		dnsFilter := bson.M{
			"subdomain": result.Domain,
			"domain":    currentSubdomain.Domain,
		}
		var lastDNSRecord models.DNS
		err := dnsColl.FindOne(context.Background(), dnsFilter, options.FindOne().SetSort(bson.M{"resolution_date": -1})).Decode(&lastDNSRecord)

		// Prepare new DNS record
		newDNSRecord := models.DNS{
			ResolutionDate: bson.NewDateTimeFromTime(now),
			Domain:         currentSubdomain.Domain,
			Subdomain:      result.Domain,
			CnameRecords:   result.Records["cname"],
			ARecords:       result.Records["a"],
			AAAARecords:    result.Records["aaaa"],
			NSRecords:      result.Records["ns"],
			PTRRecords:     result.Records["ptr"],
			MXRecords:      result.Records["mx"],
			TXTRecords:     result.Records["txt"],
		}

		// Check if we have A or AAAA records
		hasIPRecords := len(result.Records["a"]) > 0 || len(result.Records["aaaa"]) > 0

		// Check if we have any records at all
		hasAnyRecords := false
		for _, records := range result.Records {
			if len(records) > 0 {
				hasAnyRecords = true
				break
			}
		}

		var newStatus models.StatusType
		if err == mongo.ErrNoDocuments {
			// No previous DNS data
			if hasIPRecords {
				newStatus = models.FreshResolved
			} else {
				newStatus = models.UnresolvedSubdomain
			}
		} else {
			// Has previous DNS data
			if !hasIPRecords {
				newStatus = models.LastResolved
			} else if currentSubdomain.DNSStatus == models.LastResolved {
				newStatus = models.FreshResolved
			}
		}

		// Update subdomain status if needed
		if newStatus != "" {
			_, err = subdomainsColl.UpdateOne(
				context.Background(),
				bson.M{"domain": currentSubdomain.Domain, "name": currentSubdomain.Name},
				bson.M{"$set": bson.M{
					"dns_status": newStatus,
					"updated_at": bson.NewDateTimeFromTime(now),
				}},
			)
			if err != nil {
				log.Printf("failed to update subdomain status for %s: %v", currentSubdomain.Name, err)
			}
		}

		// Insert new DNS record if we have any records
		if hasAnyRecords {
			_, err = dnsColl.InsertOne(context.Background(), newDNSRecord)
			if err != nil {
				log.Printf("failed to insert DNS record for %s: %v", currentSubdomain.Name, err)
			}
		}
	}

	return nil
}

func httpxTask() error {
	log.Println("Running httpx task")
	// Httpx job

	return nil
}
