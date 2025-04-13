package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
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
		payload, err := model.NewTestPayload(&test)
		if err != nil {
			return err
		}
		testPayloads[i] = payload
	}
	return c.Status(fiber.StatusOK).JSON(&testPayloads)
}

func AddTest(c *fiber.Ctx) error {
	payload := &model.TestWritePayload{}
	err := util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	test := model.NewTest(payload)
	result := db.GetDB().Create(&test)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return service.SendTestResponse(c, &test, fiber.StatusCreated)
}

func GetTest(c *fiber.Ctx) error {
	testIDs, err := util.ReadIdsFromParams(c, []string{"id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}
	test, err := service.GetTestByID(testIDs["id"], true)
	if err != nil {
		return err
	}

	return service.SendTestResponse(c, test, fiber.StatusOK)
}

func DeleteTest(c *fiber.Ctx) error {
	testIDs, err := util.ReadIdsFromParams(c, []string{"id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}
	test, err := service.GetTestByID(testIDs["id"], false)
	if err != nil {
		return err
	}

	result := db.GetDB().Delete(test)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func UpdateTest(c *fiber.Ctx) error {
	payload := &model.TestWritePayload{}
	err := util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	testIDs, err := util.ReadIdsFromParams(c, []string{"id"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	test, err := service.GetTestByID(testIDs["id"], true)
	if err != nil {
		return err
	}

	test.UpdateFromPayload(payload)

	result := db.GetDB().Save(test)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return service.SendTestResponse(c, test, fiber.StatusOK)
}
