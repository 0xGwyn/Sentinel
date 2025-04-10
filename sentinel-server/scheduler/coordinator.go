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

	var lastJob struct {
		StartTime time.Time `bson:"start_time"`
		EndTime   time.Time `bson:"end_time"`
	}

	if err := result.Decode(&lastJob); err != nil {
		return true
	}

	// If the last job hasn't finished, don't start a new one
	if lastJob.EndTime.IsZero() {
		return false
	}

	return true
}

func (c *Coordinator) StartJob(jobType JobType) error {
	_, err := c.collection.InsertOne(
		context.Background(),
		bson.M{
			"type":       jobType,
			"start_time": time.Now(),
		},
	)
	return err
}

func (c *Coordinator) EndJob(jobType JobType) error {
	_, err := c.collection.UpdateOne(
		context.Background(),
		bson.M{
			"type":     jobType,
			"end_time": bson.M{"$exists": false},
		},
		bson.M{
			"$set": bson.M{"end_time": time.Now()},
		},
	)
	return err
}
