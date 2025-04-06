package database

import (
	"context"
	"errors"

	"github.com/0xgwyn/sentinel/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func InitDB() error {
	uri, err := config.LoadEnv("MONGODB_URI")
	if err != nil {
		return errors.New("you must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	dbName, err := config.LoadEnv("DATABASE")
	if err != nil {
		return err
	}
	db = client.Database(dbName)

	return nil
}

func InitIndexes() error {
	ctx := context.Background()

	// Subdomains
	_, err := GetDBCollection("subdomains").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "domain", Value: 1}, {Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "watch", Value: 1}}},
	})
	if err != nil {
		return err
	}

	// Domains
	_, err = GetDBCollection("domains").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// DNS
	_, err = GetDBCollection("dns").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "domain", Value: 1}, {Key: "subdomain", Value: 1}, {Key: "resolution_date", Value: -1}},
	})
	if err != nil {
		return err
	}

	// HTTP
	_, err = GetDBCollection("http").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "domain", Value: 1}, {Key: "subdomain", Value: 1}, {Key: "scanning_date", Value: -1}},
	})
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
