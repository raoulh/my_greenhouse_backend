package app

import (
	"fmt"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"
	"github.com/gofiber/fiber/v2"
)

func (a *AppServer) apiNotifGet(c *fiber.Ctx, notifType uint) (err error) {
	user := c.Locals("user")
	u, ok := user.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Errorf("unauthorized"),
		})
	}

	settings, err := models.GetNotifSettings(u, notifType)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(settings)
}

func (a *AppServer) apiNotifSet(c *fiber.Ctx, notifType uint) (err error) {
	user := c.Locals("user")
	u, ok := user.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Errorf("unauthorized"),
		})
	}

	n := new(models.NotifSettings)
	if err := c.BodyParser(n); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	err = models.SetNotifSettings(u, n)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "settings changed",
	})
}

type NotifHwId struct {
	Token  string `json:"token" xml:"token" form:"token"`
	HwType uint   `json:"hw" xml:"hw" form:"hw"`
}

func (a *AppServer) apiNotifId(c *fiber.Ctx) (err error) {
	user := c.Locals("user")
	u, ok := user.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Errorf("unauthorized"),
		})
	}

	n := new(NotifHwId)
	if err := c.BodyParser(n); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	err = models.UpdateNotifToken(u, n.Token, n.HwType)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "settings changed",
	})
}
