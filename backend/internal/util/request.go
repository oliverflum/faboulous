package util

import (
	"github.com/gofiber/fiber/v2"
)

func ReadIdsFromParams(c *fiber.Ctx, paramNames []string) (map[string]uint, error) {
	ids := make(map[string]uint)
	for _, paramName := range paramNames {
		id, err := c.ParamsInt(paramName)
		if err != nil {
			return nil, err
		}
		ids[paramName] = uint(id)
	}
	return ids, nil
}
