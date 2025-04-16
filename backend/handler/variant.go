package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
	"github.com/oliverflum/faboulous/backend/util"
)

func ListVariants(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	variants := make([]model.Variant, 0)
	result := db.GetDB().Preload("Features").Where("test_id = ?", ids["testId"]).Find(&variants)
	if result.Error != nil {
		return util.HandleGormError(result)
	} else if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No variants found for this test")
	}

	variantPayloads := make([]model.VariantPayload, len(variants))
	for i, variant := range variants {
		payload, err := service.NewVariantPayload(variant, db.GetDB())
		if err != nil {
			return err
		}

		variantPayloads[i] = payload
	}

	return c.Status(fiber.StatusOK).JSON(variantPayloads)
}

func AddVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	payload := &model.VariantWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return c.Status(valErr.Code).SendString(valErr.Message)
	}

	// Check if test exists
	test, err := service.FindTestByID(ids["testId"], false)
	if err != nil {
		return err
	}

	if service.CheckVariantExists(payload.Name, ids["testId"]) {
		return fiber.NewError(fiber.StatusBadRequest, "Variant with this name already exists for this test")
	}

	if err := service.CheckVariantSizeConstraints(db.GetDB(), test, nil, payload.Size); err != nil {
		return err
	}

	variant := service.NewVariant(*payload)
	variant.TestID = test.ID
	result := db.GetDB().Create(&variant)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return service.SendVariantResponse(c, variant, fiber.StatusCreated)
}

func UpdateVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	payload := &model.VariantWritePayload{}
	valErr := util.ParseAndValidatePayload(c, payload)
	if valErr != nil {
		return c.Status(valErr.Code).SendString(valErr.Message)
	}

	variant, err := service.FindVariantById(ids["variantId"], false)
	if err != nil {
		return err
	}

	test, err := service.FindTestByID(ids["testId"], false)
	if err != nil {
		return err
	}

	if err := service.CheckVariantSizeConstraints(db.GetDB(), test, variant, payload.Size); err != nil {
		return err
	}

	variant.UpdateFromPayload(*payload)

	result := db.GetDB().Save(variant)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return service.SendVariantResponse(c, *variant, fiber.StatusOK)
}

func DeleteVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	variant, err := service.FindVariantById(ids["variantId"], false)
	if err != nil {
		return err
	}

	result := db.GetDB().Unscoped().Delete(variant)
	if result.Error != nil {
		return util.HandleGormError(result)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
