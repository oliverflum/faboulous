package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/internal/util"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/service"
)

func ListVariants(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	variants := make([]model.Variant, 0)
	result := db.GetDB().Preload("Features").Where("test_id = ?", ids["testId"]).Find(&variants)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	variantPayloads := make([]model.VariantPayload, len(variants))
	for i, variant := range variants {
		payload, err := model.NewVariantPayload(variant, db.GetDB())
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
	err = util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	// Check if test exists
	test, err := service.GetTestByID(ids["testId"], true)
	if err != nil {
		return err
	}

	if err := service.CheckVariantExists(payload.Name, ids["testId"]); err != nil {
		return err
	}

	if err := service.CheckVariantSizeConstraints(db.GetDB(), test, nil, payload.Size); err != nil {
		return err
	}

	variant := model.NewVariant(*payload)
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
	err = util.ValidatePayload(c, payload)
	if err != nil {
		return err
	}

	variant, err := service.GetVariant(ids["variantId"], false)
	if err != nil {
		return err
	}

	test, err := service.GetTestByID(ids["testId"], true)
	if err != nil {
		return err
	}

	if err := service.CheckVariantSizeConstraints(db.GetDB(), test, variant, payload.Size); err != nil {
		return err
	}

	// Update variant
	if err := variant.UpdateFromPayload(*payload); err != nil {
		return util.SendErrorRes(c, err)
	}

	result := db.GetDB().Save(variant)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return service.SendVariantResponse(c, *variant, fiber.StatusOK)
}

func DeleteVariant(c *fiber.Ctx) error {
	ids, err := util.ReadIdsFromParams(c, []string{"testId", "variantId"})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid test ID")
	}

	variant, err := service.GetVariant(ids["variantId"], false)
	if err != nil {
		return err
	}

	result := db.GetDB().Delete(variant)
	if result.Error != nil {
		return util.SendErrorRes(c, util.HandleGormError(result.Error))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
