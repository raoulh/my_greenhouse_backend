package app

import (
	"strings"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"
	"github.com/gofiber/fiber/v2"
)

func NewTokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		var deviceID string

		//get token
		headerValue := c.Get("authorization")
		if headerValue != "" {
			split := strings.SplitN(headerValue, " ", 2)

			if len(split) == 2 && strings.ToLower(split[0]) == "bearer" {
				token = split[1]
			}
		}
		c.Locals("token", token)

		//get device ID
		headerValue = c.Get("x-device-id")
		if headerValue != "" {
			deviceID = headerValue
		}
		c.Locals("device_id", deviceID)

		//get User from models
		if token != "" && deviceID != "" {
			u, err := models.GetUserByTokenAndID(token, deviceID, false)
			if err == nil {
				c.Locals("user", u)
			}
		}

		return c.Next()
	}
}
