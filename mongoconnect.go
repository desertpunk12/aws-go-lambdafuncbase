package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

var DB *mongo.Client

func MongoConnect() *mongo.Client {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	connectAndPing := func() error {
		err = DB.Connect(ctx)
		if err != nil {
			log.Fatalf("Failed to connect to cluster: %v", err)
			return err
		}
		// Force a connection to verify our connection string
		err = DB.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to ping cluster: %v", err)
			return err
		}
		return nil
	}

	// Return and reuse old connections
	if DB != nil {
		if DB.Ping(ctx, nil) == nil {
			fmt.Println("Reusing existing connection")
			return DB
		}
		err = connectAndPing()
		if err != nil {
			log.Fatalf("Failed to connect or ping cluster: %v", err)
			return nil
		}
		fmt.Println("Connected to MongoDB and successfully pinged!")

		return DB
	}

	mongoURI := fmt.Sprintf(ConnectionStringTemplate, Username, Password, ClusterEndpoint, DBName)
	// Create a new client and connect to the server
	DB, err = mongo.NewClient(options.Client().ApplyURI(mongoURI).SetRetryWrites(false))
	if err != nil {
		log.Fatalf("DocDb Error: Failed to create client: %v", err)
		return nil
	}

	err = connectAndPing()
	if err != nil {
		log.Fatalf("Failed to connect or ping cluster: %v", err)
		return nil
	}

	return DB

}

func Test(t *testing.T) error {
	colName := "colName"
	collection := MongoConnect().Database(DBName).Collection(colName)

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	res := collection.FindOne(ctx, nil)
	if res.Err() != nil {
		t.Errorf("Error: %v", res.Err())
		return res.Err()
	}

	return nil
}
