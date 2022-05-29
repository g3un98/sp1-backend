package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type user struct {
	AppId    string `json:"app_id" bson:"app_id"`
	AppPw    string `json:"app_pw" bson:"app_pw"`
	AppEmail string `json:"app_email,omitempty" bson:"app_email,omitempty"`
}

type member struct {
	AppId   string `json:"app_ip" bson:"app_ip"`
	IsAdmin int    `json:"is_admin" bson:"is_admin"`
}

type group struct {
	group_id   primitive.ObjectID `json:"idx" bson:"_id,omitempty"`
	Ott        string            `json:"ott" bson:"ott"`
    Account    account           `json:"account" bson:"account"`
    Updatetime int64             `json:"updatetime" bson:"updatetime"`
	Members    []member          `json:"members" bson:"members"`
}

const (
	DB_URI  = "mongodb://db:27017"
	DB_NAME = "ott"
)

func connectDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)

	clientOptions := options.Client().ApplyURI(DB_URI).SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	})
	client, err := mongo.Connect(ctx, clientOptions)
    checkError(err)

	err = client.Ping(ctx, readpref.Primary())
    checkError(err)

	return client, ctx, cancel, nil
}

func getCollection(client *mongo.Client, colName string) *mongo.Collection {
	return client.Database(DB_NAME).Collection(colName)
}

func addUser(c *fiber.Ctx) error {
	client, ctx, cancel, err := connectDB()
    checkError(err)
	defer cancel()
	defer client.Disconnect(ctx)

	var user user
	err = c.BodyParser(&user)
    checkError(err)

	filter := bson.M{"app_id": user.AppId}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
    checkError(err)

	if num == 0 {
		_, err = getCollection(client, "user").InsertOne(ctx, user)
        checkError(err)

		return c.SendStatus(fiber.StatusCreated)
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}

func delUser(c *fiber.Ctx) error {
	client, ctx, cancel, err := connectDB()
    checkError(err)
	defer cancel()
	defer client.Disconnect(ctx)

	var user user
	err = c.BodyParser(&user)
    checkError(err)

	filter := bson.M{"app_id": user.AppId, "app_pw": user.AppPw}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
    checkError(err)

	if num == 1 {
		_, err = getCollection(client, "user").DeleteOne(ctx, user)
        checkError(err)

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}

func setUser(c *fiber.Ctx) error {
	client, ctx, cancel, err := connectDB()
    checkError(err)
	defer cancel()
	defer client.Disconnect(ctx)

	var user user
	err = c.BodyParser(&user)
    checkError(err)

	filter := bson.M{"app_id": user.AppId, "app_pw": user.AppPw}
	update := bson.M{"$set": bson.M{"app_id": user.AppId, "app_pw": user.AppPw, "app_email": user.AppEmail}}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
	checkError(err)

	if num == 1 {
		_, err = getCollection(client, "user").UpdateOne(ctx, filter, update)
        checkError(err)

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}

func login(c *fiber.Ctx) error {
	client, ctx, cancel, err := connectDB()
	checkError(err)
	defer cancel()
	defer client.Disconnect(ctx)

	var user user
	err = c.BodyParser(&user)
    checkError(err)

	filter := bson.M{"app_id": user.AppId, "app_pw": user.AppPw}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
	checkError(err)

	if num == 1 {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func getGroup(c *fiber.Ctx) error {
	client, ctx, cancel, err := connectDB()
	checkError(err)
	defer cancel()
	defer client.Disconnect(ctx)

	filter := bson.M{"idx": c.Params("idx")}

	num, err := getCollection(client, "group").CountDocuments(ctx, filter)
	checkError(err)

	var group bson.M
	if num == 1 {
		err = getCollection(client, "group").FindOne(ctx, filter).Decode(&group)
        checkError(err)
		body, err := bson.Marshal(group)
        checkError(err)
		return c.Send(body)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func addGroup(c *fiber.Ctx) error {
	client, ctx, cancel, err := connectDB()
    checkError(err)
	defer cancel()
	defer client.Disconnect(ctx)

	var parser struct {
		AppId string `json:"app_id" bson:"app_id"`
		Ott   string `json:"ott" bson:"ott"`
		OttId string `json:"ott_id" bson:"ott_id"`
		OttPw string `json:"ott_pw" bson:"ott_pw"`
	}
	err = c.BodyParser(&parser)
    checkError(err)

    filter := bson.M{ "ott": parser.Ott, "account.id": parser.OttId, "account.pw": parser.OttPw }
	num, err := getCollection(client, "group").CountDocuments(ctx, filter)
	checkError(err)

	var group group
	switch num {
	case 0:
        account, err := getAccount(parser.Ott, parser.OttId, parser.OttPw)
        if err != nil {
            return err
        }
        
		group.Ott = parser.Ott
		group.Account = *account
		group.Updatetime = time.Now().Unix()
		group.Members = []member{{
			AppId:   parser.AppId,
			IsAdmin: 1,
		}}

		_, err = getCollection(client, "group").InsertOne(ctx, group)
        checkError(err)

		return c.SendStatus(fiber.StatusOK)
	case 1:
        filter2 := bson.M{ "ott": parser.Ott, "account.id": parser.OttId, "account.pw": parser.OttPw, "members.app_id": parser.AppId }
	    num, err := getCollection(client, "group").CountDocuments(ctx, filter2)
	    checkError(err)
        if num == 1 {
	        return c.SendStatus(fiber.StatusUnauthorized)
        }

        update := bson.M{"$push": bson.M{ "members": member{parser.AppId, 0} }, "$set": bson.M{ "updatetime": time.Now().Unix() }}
		_, err = getCollection(client, "group").UpdateOne(ctx, filter, update)
        checkError(err)

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusBadRequest)
}
