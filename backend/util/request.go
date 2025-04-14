package util

import (
	"github.com/gofiber/fiber/v2"
)

func ReadIdsFromParams(c *fiber.Ctx, paramNames []string) (map[string]uint, *fiber.Error) {
	ids := make(map[string]uint)
	for _, paramName := range paramNames {
		id, err := c.ParamsInt(paramName)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid "+paramName+" ID")
		}
		ids[paramName] = uint(id)
	}
	return ids, nil
}
