package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/db"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/util"
	"gorm.io/gorm"
)

func SendTestResponse(c *fiber.Ctx, test *model.Test, statusCode int) error {
	payload, err := model.NewTestPayload(test)
	if err != nil {
		return err
	}
	return c.Status(statusCode).JSON(payload)
}

func GetTestByID(id uint, preloadVariants bool) (*model.Test, *fiber.Error) {
	var test model.Test
	var result *gorm.DB
	if preloadVariants {
		result = db.GetDB().
			Preload("Variants").
			First(&test, id)
	} else {
		result = db.GetDB().First(&test, id)
	}
	if result.Error != nil {
		return nil, util.HandleGormError(result)
	} else if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Test not found")
	}
	return &test, nil
}

func GetAllTests(preloadVariants bool) ([]model.Test, *fiber.Error) {
	var tests []model.Test
	var result *gorm.DB
	if preloadVariants {
		result = db.GetDB().Preload("Variants").Find(&tests)
	} else {
		result = db.GetDB().Find(&tests)
	}
	if result.Error != nil {
		return nil, util.HandleGormError(result)
	}
	return tests, nil
}

func CheckIfMethodAllowed(testMethod string) *fiber.Error {
	if testMethod == model.HASH {
		var count int64
		result := db.GetDB().Model(&model.Test{}).Where("method = ?", model.HASH).Count(&count)
		if result.Error != nil {
			return util.HandleGormError(result)
		}
		exists := count > 0
		if exists {
			return fiber.NewError(fiber.StatusConflict, "Only one test with HASH method is allowed")
		}
		return nil
	}
	return nil
}
