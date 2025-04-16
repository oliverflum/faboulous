package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
	"github.com/oliverflum/faboulous/backend/util"
)

func ListTests(c *fiber.Ctx) error {
	tests, err := service.GetAllTests(true)
	if err != nil {
		return err
	}

	if len(tests) == 0 {
		return c.Status(fiber.StatusOK).JSON(&tests)
	}

	testPayloads := make([]model.TestPayload, len(tests))
	for i, test := range tests {
		payload, err := service.NewTestPayload(test)
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

	err := service.CheckIfMethodAllowed(payload.Method)
	if err != nil {
		return err
	}

	test := service.NewTest(payload)
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
	test, err := service.FindTestByID(ids["testId"], true)
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
	test, err := service.FindTestByID(ids["testId"], false)
	if err != nil {
		return err
	}

	result := db.GetDB().Unscoped().Delete(test)
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

	err = service.CheckIfMethodAllowed(payload.Method)
	if err != nil {
		return err
	}

	test, err := service.FindTestByID(ids["testId"], false)
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
