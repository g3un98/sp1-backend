package main

import (
    "net/http"

	"github.com/gofiber/fiber/v2"
)

const DB_API_BASE_URL = "http://api:12390"

func addUser(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodPost, DB_API_BASE_URL + "/users", c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}

func delUser(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodDelete, DB_API_BASE_URL + "/users", c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}

func setUser(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodPut, DB_API_BASE_URL + "/users", c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}

func login(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodPost, DB_API_BASE_URL + "/login", c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}

func addGroup(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodPost, DB_API_BASE_URL + "/group", c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}

func delGroup(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodDelete, DB_API_BASE_URL + "/group", c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}

func getGroup(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodGet, DB_API_BASE_URL + "/otts/group/" + c.Params("idx"), c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    if res.StatusCode == 200 {
        return c.SendStream(res.Body)
    }

    return c.SendStatus(res.StatusCode)
}

func setGroup(c *fiber.Ctx) error {
    client := &http.Client{}

    req, err := http.NewRequest(http.MethodGet, DB_API_BASE_URL + "/otts/group/" + c.Params("idx"), c.Context().RequestBodyStream())
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    defer res.Body.Close()

    return c.SendStatus(res.StatusCode)
}
