package controllers

import (
	"encoding/json"
	"omokogo/globals"
	"omokogo/models"
	"omokogo/queries"

	"github.com/gofiber/fiber/v2"
)

func LoginController(c *fiber.Ctx) error {
	sess, _ := globals.Store.Get(c)

	body := models.LoginReq{}
	_ = json.Unmarshal(c.Body(), &body)

	user, err := queries.LoginUser(body)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{"error": true, "message": "wrong username or wrong password"})
	}

	sess.Set("id", user.Id)
	sess.Save()

	return c.Status(200).JSON(fiber.Map{})
}

func LoginViewController(c *fiber.Ctx) error {
	return c.Render("pages/auth/login", fiber.Map{
		"Title": "login",
	}, "layout")
}

func RegisterController(c *fiber.Ctx) error {
	sess, _ := globals.Store.Get(c)

	body := models.RegisterReq{}
	_ = json.Unmarshal(c.Body(), &body)

	nameExists := queries.CheckUserName(body.Name)
	if nameExists {
		return c.Status(400).JSON(
			fiber.Map{"error": true, "message": "username exists"})
	}

	emailExists := queries.CheckUserEmail(body.Email)
	if emailExists {
		return c.Status(400).JSON(
			fiber.Map{"error": true, "message": "email exists"})
	}

	user, err := queries.CreateUser(body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	sess.Set("id", user.Id)
	sess.Save()

	return c.Status(200).JSON(fiber.Map{})
}

func RegisterViewController(c *fiber.Ctx) error {
	return c.Render("pages/auth/register", fiber.Map{
		"Title": "register",
	}, "layout")
}
