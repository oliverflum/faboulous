package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
)

func ListFeatures(c *fiber.Ctx) error {
	var features []model.Feature

	util.GetDB().Find(&features)

	return c.Status(200).JSON(&features)
}

func AddFeature(c *fiber.Ctx) error {
	feature := new(model.Feature)

	if err := c.BodyParser(feature); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	errors := util.ValidateStruct(feature)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	result := util.GetDB().Create(&feature)

	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(500).Send(nil)
	}
	return c.Status(200).JSON(feature)
}

func GetFeature(c *fiber.Ctx) error {
	id := c.Params("id")
	var feature model.Feature

	result := util.GetDB().Find(&feature, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&feature)
}

func DeleteFeature(c *fiber.Ctx) error {
	id := c.Params("id")

	result := util.GetDB().Delete(&model.Feature{}, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(204).Send(nil)
}
