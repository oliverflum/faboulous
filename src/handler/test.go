package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
)

func ListTests(c *fiber.Ctx) error {
	var tests []model.TestEntity

	util.GetDB().Find(&tests)

	if len(tests) == 0 {
		return c.Status(200).JSON(&tests)
	}

	return c.Status(200).JSON(&tests)
}

func AddTest(c *fiber.Ctx) error {
	var payload model.TestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	test := model.NewTestEntity(payload)
	result := util.GetDB().Create(&test)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(500).Send(nil)
	}

	return c.Status(200).JSON(model.NewTestPayload(test))
}

func GetTest(c *fiber.Ctx) error {
	id := c.Params("id")
	var test model.TestEntity

	result := util.GetDB().Find(&test, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&test)
}

func DeleteTest(c *fiber.Ctx) error {
	id := c.Params("id")

	result := util.GetDB().Delete(&model.TestEntity{}, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(204).Send(nil)
}

func UpdateTest(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload model.TestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	var existingTest model.TestEntity
	result := util.GetDB().First(&existingTest, id)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	test := model.NewTestEntity(payload)
	test.ID = existingTest.ID
	saveResult := util.GetDB().Save(&test)
	if saveResult.Error != nil {
		return c.Status(500).SendString(saveResult.Error.Error())
	}

	return c.Status(200).JSON(model.NewTestPayload(test))
}
