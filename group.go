package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type member struct {
	AppId   string `json:"app_id" bson:"app_id"`
	IsAdmin int    `json:"is_admin" bson:"is_admin"`
}

type group struct {
	GroupId    *primitive.ObjectID   `json:"groupId" bson:"_id,omitempty"`
	Ott        string   `json:"ott" bson:"ott"`
	Account    account  `json:"account" bson:"account"`
	UpdateTime int64    `json:"update_time" bson:"update_time"`
	Members    []member `json:"members" bson:"members"`
}

func getGroup(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

	_id, err := primitive.ObjectIDFromHex(c.Params("groupId"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	filter := bson.M{"_id": _id}

	num, err := getCollection(client, "group").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var group bson.M
	if num == 1 {
		if err = getCollection(client, "group").FindOne(ctx, filter).Decode(&group); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		body, err := bson.Marshal(group)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.Send(body)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func addGroup(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

	var parser struct {
		AppId string `json:"app_id" bson:"app_id"`
		Ott   string `json:"ott" bson:"ott"`
		OttId string `json:"ott_id" bson:"ott_id"`
		OttPw string `json:"ott_pw" bson:"ott_pw"`
	}

	if err = c.BodyParser(&parser); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	filter := bson.M{"ott": parser.Ott, "account.id": parser.OttId, "account.pw": parser.OttPw}
	num, err := getCollection(client, "group").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var group group
	switch num {
	case 0:
		account, err := getAccount(parser.Ott, parser.OttId, parser.OttPw)
		if err != nil {
			return err
		}

		group.Ott = parser.Ott
		group.Account = *account
		group.UpdateTime = time.Now().Unix()
		group.Members = []member{{
			AppId:   parser.AppId,
			IsAdmin: 1,
		}}

		if _, err = getCollection(client, "group").InsertOne(ctx, group); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	case 1:
		filter2 := bson.M{"ott": parser.Ott, "account.id": parser.OttId, "account.pw": parser.OttPw, "members.app_id": parser.AppId}
		num, err := getCollection(client, "group").CountDocuments(ctx, filter2)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if num == 1 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		filter3 := bson.M{"app_id": parser.AppId}
		num, err = getCollection(client, "user").CountDocuments(ctx, filter3)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if num != 1 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		update := bson.M{"$push": bson.M{"members": member{parser.AppId, 0}}, "$set": bson.M{"update_time": time.Now().Unix()}}
		if _, err = getCollection(client, "group").UpdateOne(ctx, filter, update); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusBadRequest)
}

func delGroup(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

	var parser struct {
		AppId string `json:"app_id" bson:"app_id"`
	}
	if err = c.BodyParser(&parser); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_id, err := primitive.ObjectIDFromHex(c.Params("groupId"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var group group
	filter := bson.M{"_id": _id, "members.app_id": parser.AppId}
	if err = getCollection(client, "group").FindOne(ctx, filter).Decode(&group); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if containAdminMembers(group.Members, parser.AppId) {
		if _, err = getCollection(client, "group").DeleteOne(ctx, filter); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func setGroup(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

	var parser struct {
		OttPw      string     `json:"ott_pw" bson:"ott_pw"`
		Payment    payment    `json:"payment" bson:"payment"`
		Membership membership `json:"membership" bson:"membership"`
	}
	if err = c.BodyParser(&parser); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	_id, err := primitive.ObjectIDFromHex(c.Params("groupId"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	filter := bson.M{"_id": _id}
	num, err := getCollection(client, "group").CountDocuments(ctx, filter)

	if num == 1 {
		update := bson.M{"$set": bson.M{
			"account.pw":              parser.OttPw,
			"account.payment.type":    parser.Payment.Type,
			"account.payment.detail":  parser.Payment.Detail,
			"account.payment.next":    parser.Payment.Next,
			"account.membership.type": parser.Membership.Type,
			"account.membership.cost": parser.Membership.Cost,
		}}
		if _, err = getCollection(client, "group").UpdateOne(ctx, filter, update); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusNotFound)
}
