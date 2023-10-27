package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func IndexController(c *fiber.Ctx) error {
	return c.Render("pages/index", fiber.Map{
		"Title": "hi",
	}, "layout")
}

func LoginController(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)

	body := LoginReq{}
	_ = json.Unmarshal(c.Body(), &body)

	user, err := LoginUser(body)
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
	sess, _ := Store.Get(c)

	body := RegisterReq{}
	_ = json.Unmarshal(c.Body(), &body)

	nameExists := CheckUserName(body.Name)
	if nameExists {
		return c.Status(400).JSON(
			fiber.Map{"error": true, "message": "username exists"})
	}

	emailExists := CheckUserEmail(body.Email)
	if emailExists {
		return c.Status(400).JSON(
			fiber.Map{"error": true, "message": "email exists"})
	}

	user, err := CreateUser(body)
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

func MeController(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	id, ok := sess.Get("id").(string)
	if ok != true {
		return c.Status(400).SendString("login please")
	}
	return c.SendString(id)
}

func TestController(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	_, ok := sess.Get("id").(string)
	if ok != true {
		return c.Status(400).SendString("login please")
	}
	return c.Render("pages/test", fiber.Map{
		"Title": "register",
	}, "layout")
}

