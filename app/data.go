package app

import (
	"fmt"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"
	"github.com/gofiber/fiber/v2"
)

func (a *AppServer) apiGetDataFull(c *fiber.Ctx) (err error) {
	user := c.Locals("user")
	u, ok := user.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Errorf("unauthorized"),
		})
	}

	fullUser, err := models.GetFullUser(u.ID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fullUser)
}
