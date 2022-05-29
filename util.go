package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func containMembers(s []member, v string) bool {
	for _, vv := range s {
		if v == vv.AppId {
			return true
		}
	}
	return false
}

func checkError(err error) {
    if err != nil {
        log.Fatal(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
    }
}
