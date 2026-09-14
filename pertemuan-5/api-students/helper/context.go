package helper

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
)

// LocalsAuthUser adalah kunci penyimpanan identitas di context request.
const LocalsAuthUser = "authUser"

func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	user, ok := c.Locals(LocalsAuthUser).(model.AuthUser)
	return user, ok
}