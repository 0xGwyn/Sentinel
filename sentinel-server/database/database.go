package database

import (
	"context"
	"errors"

	"github.com/0xgwyn/sentinel/common"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func InitDB() error {
	uri, err := common.LoadEnv("MONGODB_URI")
	if err != nil {
		return errors.New("you must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	dbName, err := common.LoadEnv("DATABASE")
	if err != nil {
		return err
	}
	db = client.Database(dbName)

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
