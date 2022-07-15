package app

import (
	"fmt"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"
	"github.com/gofiber/fiber/v2"
)

type AuthUser struct {
	Name     string `json:"username" xml:"username" form:"username"`
	Pass     string `json:"pass" xml:"pass" form:"pass"`
	DeviceID string `json:"device_id" xml:"device_id" form:"device_id"`
}

func (a *AppServer) apiLogin(c *fiber.Ctx) (err error) {
	u := new(AuthUser)

	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if u.Name == "" || u.Pass == "" || u.DeviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Errorf("bad arguments"),
		})
	}

	at, err := models.Login(u.Name, u.Pass, u.DeviceID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(at)
}

func (a *AppServer) apiCheckToken(c *fiber.Ctx) (err error) {
	user := c.Locals("user")
	u, ok := user.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Errorf("unauthorized"),
		})
	}

	err = models.CheckToken(u)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	models.UpdateLastLogin(u)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "token valid",
	})
}
