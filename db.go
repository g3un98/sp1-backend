package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type user struct {
    AppId string `json:"app_id" bson:"app_id"`
    AppPw string `json:"app_pw" bson:"app_pw"`
    AppEmail string `json:"app_email,omitempty" bson:"app_email,omitempty"`
}

type group struct {
    Idx int `json:"idx" bson:"idx"`
    Ott string `json:"ott" bson:"ott"`
    Account account `json:"account" bson:"account"`
    Updatetime int64 `json:"updatetime" bson:"updatetime"`
    Members []struct{} `json:"members" bson:"members"`
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

    var user user
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

    filter := bson.M{ "app_id": user.AppId }

    num, err := getCollection(client, "user").CountDocuments(ctx, filter)
    if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
    }

    if num == 0 {
        if _, err := getCollection(client, "user").InsertOne(ctx, user); err != nil { 
		    return c.SendStatus(fiber.StatusBadRequest)
        }

        return c.SendStatus(fiber.StatusCreated)
    }

    return c.SendStatus(fiber.StatusUnauthorized)
}

func delUser(c *fiber.Ctx) error {
    client, ctx, cancel, err := connectDB()
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer cancel()
    defer client.Disconnect(ctx)

    var user user
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

    filter := bson.M{ "app_id": user.AppId, "app_pw": user.AppPw }

    num, err := getCollection(client, "user").CountDocuments(ctx, filter)
    if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
    }

    if num == 1 {
        if _, err := getCollection(client, "user").DeleteOne(ctx, user); err != nil { 
		    return c.SendStatus(fiber.StatusBadRequest)
        }

        return c.SendStatus(fiber.StatusOK)
    }

    return c.SendStatus(fiber.StatusUnauthorized)
}

func setUser(c *fiber.Ctx) error {
    client, ctx, cancel, err := connectDB()
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer cancel()
    defer client.Disconnect(ctx)

    var user user
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

    filter := bson.M{ "app_id": user.AppId, "app_pw": user.AppPw }
    update := bson.M{ "$set": bson.M{ "app_id": user.AppId, "app_pw": user.AppPw, "app_email": user.AppEmail} }

    num, err := getCollection(client, "user").CountDocuments(ctx, filter)
    if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
    }

    if num == 1 {
        if _, err := getCollection(client, "user").UpdateOne(ctx, filter, update); err != nil { 
		    return c.SendStatus(fiber.StatusBadRequest)
        }

        return c.SendStatus(fiber.StatusOK)
    }

    return c.SendStatus(fiber.StatusUnauthorized)
}

func login(c *fiber.Ctx) error {
    client, ctx, cancel, err := connectDB()
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer cancel()
    defer client.Disconnect(ctx)

    var user user
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

    filter := bson.M{ "app_id": user.AppId, "app_pw": user.AppPw }

    num, err := getCollection(client, "user").CountDocuments(ctx, filter)
    if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
    }

    if num == 1 {
        return c.SendStatus(fiber.StatusOK)
    }

    return c.SendStatus(fiber.StatusNotFound)
}

func getGroup(c *fiber.Ctx) error {
    client, ctx, cancel, err := connectDB()
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer cancel()
    defer client.Disconnect(ctx)

    filter := bson.M{ "idx": c.Params("idx") }

    num, err := getCollection(client, "group").CountDocuments(ctx, filter)
    if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
    }

    var group bson.M
    if num == 1 {
        if err := getCollection(client, "group").FindOne(ctx, filter).Decode(&group); err != nil {
		    return c.SendStatus(fiber.StatusBadRequest)
        }
        body, err := bson.Marshal(group)
        if err != nil {
		    return c.SendStatus(fiber.StatusBadRequest)
        }
        return c.Send(body)
    }

    return c.SendStatus(fiber.StatusNotFound)
}
