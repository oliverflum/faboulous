package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/db"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
)

func ListTests(c *fiber.Ctx) error {
	var tests []model.Test

	result := db.GetDB().
		Preload("Variants").
		Preload("Variants.Features").
		Find(&tests)
	if result.Error != nil {
		return util.SendErrorRes(c, util.HandleGormError(result.Error))
	}

	if len(tests) == 0 {
		return c.Status(fiber.StatusOK).JSON(&tests)
	}

	testPayloads := make([]model.TestPayload, len(tests))
	for i, test := range tests {
		testPayloads[i] = model.NewTestPayload(test)
	}
	return c.Status(fiber.StatusOK).JSON(&testPayloads)
}

func appendVariants(payload *model.TestPayload, test *model.Test) error {
	// Iterate over variants and check if they exist ind db
	for _, variant := range payload.Variants {
		var existingVariant model.Variant
		result := db.GetDB().First(&existingVariant, variant.Id)
		if result.RowsAffected == 0 {
			return errors.New("Variant with id " + fmt.Sprint(variant.Id) + " does not exist")
		}
		test.Variants = append(test.Variants, existingVariant)
	}
	return nil
}

func AddTest(c *fiber.Ctx) error {
	var payload model.TestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	test := model.NewTestEntity(payload)
	// Iterate over variants and check if they exist ind db
	err := appendVariants(&payload, &test)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result := db.GetDB().Create(&test)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(model.NewTestPayload(test))
}

func GetTest(c *fiber.Ctx) error {
	id := c.Params("id")
	var test model.Test

	result := db.GetDB().Find(&test, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(model.NewTestPayload(test))
}

func DeleteTest(c *fiber.Ctx) error {
	id := c.Params("id")

	result := db.GetDB().Delete(&model.Test{}, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func UpdateTest(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload model.TestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errors := util.ValidateStruct(payload)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var existingTest model.Test
	result := db.GetDB().First(&existingTest, id)
	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	test := model.NewTestEntity(payload)
	test.Variants = make([]model.Variant, len(payload.Variants))
	err := appendVariants(&payload, &test)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	test.ID = existingTest.ID
	saveResult := db.GetDB().Save(&test)
	if saveResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(saveResult.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(model.NewTestPayload(test))
}
