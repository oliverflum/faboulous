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
		Find(&tests)
	if result.Error != nil {
		return util.HandleGormError(result)
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
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	test := model.NewTest(payload)
	result := db.GetDB().Create(&test)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return service.SendTestResponse(c, &test, fiber.StatusCreated)
}

func GetTest(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return err
	}
	test, err := service.GetTestByID(ids["testId"], true)
	if err != nil {
		return err
	}

	return service.SendTestResponse(c, test, fiber.StatusOK)
}

func DeleteTest(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return err
	}
	test, err := service.GetTestByID(ids["testId"], false)
	if err != nil {
		return err
	}

	result := db.GetDB().Delete(test)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func UpdateTest(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return err
	}

	payload := &model.TestWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return valErr
	}

	test, err := service.GetTestByID(ids["testId"], true)
	if err != nil {
		return err
	}

	test.UpdateFromPayload(payload)

	result := db.GetDB().Save(test)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return service.SendTestResponse(c, test, fiber.StatusOK)
}

func Publish(c *fiber.Ctx) error {
	// Implement the logic for publishing
	return c.SendStatus(fiber.StatusOK)
}
