package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Name string `json:"name"`
}

func CreateUserController(c *fiber.Ctx) error {
	data := CreateUserRequest{}
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		return err
	}

	newUser := User{
		Id:      CreateId(),
		Name:    data.Name,
		Score:   0.0,
		Status:  UserStatusIdle,
		GameId:  "",
		Sockets: nil,
	}
	NewUser(&newUser)

	return c.Status(200).JSON(&fiber.Map{})
}
