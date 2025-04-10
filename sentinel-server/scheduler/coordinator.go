package scheduler

import (
	"context"
	"time"

	"github.com/0xgwyn/sentinel/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobType string

const (
	SubfinderJob JobType = "subfinder"
	HttpxJob     JobType = "httpx"
	DnsxJob      JobType = "dnsx"
)

type JobStatus string

const (
	JobStatusPending JobStatus = "pending"
	JobStatusSuccess JobStatus = "success"
	JobStatusFailed  JobStatus = "failed"
)

type Job struct {
	Type      JobType   `bson:"type"`
	StartTime time.Time `bson:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty"`
	Status    JobStatus `bson:"status"`
	Error     string    `bson:"error,omitempty"`
}

type Coordinator struct {
	collection *mongo.Collection
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		collection: database.GetDBCollection("jobs"),
	}
}

func (c *Coordinator) CanRun(jobType JobType) bool {
	ctx := context.Background()

	// Look for the most recent job of this type
	result := c.collection.FindOne(
		ctx,
		bson.M{"type": jobType},
		options.FindOne().SetSort(bson.M{"start_time": -1}),
	)

	if result.Err() == mongo.ErrNoDocuments {
		return true
	}

	var lastJob Job
	if err := result.Decode(&lastJob); err != nil {
		return true
	}

	// If the last job hasn't finished, don't start a new one
	return !lastJob.EndTime.IsZero()
}

func (c *Coordinator) StartJob(jobType JobType) error {
	job := Job{
		Type:      jobType,
		StartTime: time.Now(),
		Status:    JobStatusPending,
	}

	_, err := c.collection.InsertOne(context.Background(), job)
	return err
}

func (c *Coordinator) EndJob(jobType JobType, err error) error {
	status := JobStatusSuccess
	job := Job{
		EndTime: time.Now(),
		Status:  status,
	}

	if err != nil {
		job.Status = JobStatusFailed
		job.Error = err.Error()
	}

	_, updateErr := c.collection.UpdateOne(
		context.Background(),
		bson.M{
			"type":     jobType,
			"end_time": bson.M{"$exists": false},
		},
		bson.M{"$set": job},
	)
	return updateErr
}
