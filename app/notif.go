package app

import (
	"fmt"
	"strconv"

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

	prodIdStr := c.Params("prodid", "0")
	prodId, err := strconv.Atoi(prodIdStr)
	if err != nil || prodId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	settings, err := models.GetNotifSettings(u, notifType, uint(prodId))
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
	Token       string `json:"token" xml:"token" form:"token"`
	HwType      uint   `json:"hw" xml:"hw" form:"hw"`
	Locale      string `json:"locale" xml:"locale" form:"locale"`
	Development bool   `json:"dev" xml:"dev" form:"dev"`
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

	err = models.UpdateNotifToken(u, n.Token, n.HwType, n.Locale, n.Development)
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
