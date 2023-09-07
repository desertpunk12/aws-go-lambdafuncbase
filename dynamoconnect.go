package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

var dync *dynamodb.Client

func DynamoConnect() *dynamodb.Client {
	if dync != nil {
		res, err := dync.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(2),
		})
		if err == nil {
			log.Println("tables:", res.TableNames)
			return dync
		} else {
			log.Println("error:", err)
		}
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
		return nil
	}

	dync = dynamodb.NewFromConfig(cfg)
	log.Println("Created new DynamoDB Connection")
	res, err := dync.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})
	if err != nil {
		log.Println("error listing tab les:", err)
	}
	log.Println("tables:", res.TableNames)

	return dync
}
