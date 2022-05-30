package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type user struct {
	AppId    string `json:"app_id" bson:"app_id"`
	AppPw    string `json:"app_pw" bson:"app_pw"`
	AppEmail string `json:"app_email,omitempty" bson:"app_email,omitempty"`
}

func addUser(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

    var parser struct {
        AppId string `json:"app_id" bson:"app_id"`
        AppPw string `json:"app_pw" bson:"app_pw"`
    }
	if err = c.BodyParser(&parser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	filter := bson.M{"app_id": parser.AppId}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if num == 0 {
		if _, err = getCollection(client, "user").InsertOne(ctx, parser); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusCreated)
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}

func delUser(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

    var parser struct {
        AppId string `json:"app_id" bson:"app_id"`
        AppPw string `json:"app_pw" bson:"app_pw"`
    }
	if err = c.BodyParser(&parser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	filter := bson.M{"app_id": parser.AppId, "app_pw": parser.AppPw}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if num == 1 {
		if _, err = getCollection(client, "user").DeleteOne(ctx, filter); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}

func setUser(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

    var parser struct {
        AppId string `json:"app_id" bson:"app_id"`
        AppPw string `json:"app_pw" bson:"app_pw"`
        AppEmail string `json:"app_email" bson:"app_email"`
    }
	if err = c.BodyParser(&parser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	filter := bson.M{"app_id": parser.AppId, "app_pw": parser.AppPw}
	update := bson.M{"$set": bson.M{"app_id": parser.AppId, "app_pw": parser.AppPw, "app_email": parser.AppEmail}}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if num == 1 {
		if _, err = getCollection(client, "user").UpdateOne(ctx, filter, update); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}

func login(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

    var parser struct {
        AppId string `json:"app_id" bson:"app_id"`
        AppPw string `json:"app_pw" bson:"app_pw"`
    }
	if err = c.BodyParser(&parser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	filter := bson.M{"app_id": parser.AppId, "app_pw": parser.AppPw}

	num, err := getCollection(client, "user").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if num == 1 {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.StatusNotFound)
}
