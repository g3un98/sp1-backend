package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DB_URI  = "mongodb://db:27017"
	DB_NAME = "ott"
)

func newClient() (*mongo.Client, context.Context, func(), error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)

	clientOptions := options.Client().ApplyURI(DB_URI).SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return client, ctx, func() {
        defer cancel()
        defer client.Disconnect(ctx)
    }, nil
}

func getCollection(client *mongo.Client, colName string) *mongo.Collection {
	return client.Database(DB_NAME).Collection(colName)
}
