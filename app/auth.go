package app

import (
	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"
	"github.com/gofiber/fiber/v2"
)

type AuthUser struct {
	Name string `json:"username" xml:"username" form:"username"`
	Pass string `json:"pass" xml:"pass" form:"pass"`
}

func (a *AppServer) apiLogin(c *fiber.Ctx) (err error) {
	u := new(AuthUser)

	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	at, err := models.Login(u.Name, u.Pass)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(at)
}

func (a *AppServer) apiLogout(c *fiber.Ctx) (err error) {
	u := new(models.AuthToken)

	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	err = models.Logout(u)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "logged out",
	})
}
