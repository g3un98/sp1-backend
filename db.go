package main

import (
	"context"
	"time"
    "fmt"

	"github.com/gofiber/fiber/v2"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type user struct {
    appId string `json:"app_id"`
    appPw string `json:"app_pw"`
}

const (
    DB_URI = "mongodb://db:27017"
    DB_NAME = "ott"
)

func connectDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
    ctx, cancel := context.WithTimeout(context.TODO(), 1 * time.Minute)
    
    clientOptions := options.Client().ApplyURI(DB_URI).SetAuth(options.Credential{
        Username: "root",
        Password: "root",
    })
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, nil, nil, err
    }

    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        return nil, nil, nil, err
    }

    return client, ctx, cancel, nil
}

func getCollection(client *mongo.Client, colName string) *mongo.Collection {
    return client.Database(DB_NAME).Collection(colName)
}

func addUser(c *fiber.Ctx) error {
    client, ctx, cancel, err := connectDB()
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer cancel()
    defer client.Disconnect(ctx)

    var u user
	if err := c.BodyParser(&u); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
    fmt.Println(u)
    fmt.Println(string(c.Body()))

    if _, err := getCollection(client, "user").InsertOne(ctx, u); err != nil { 
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.SendStatus(fiber.StatusOK)
}
