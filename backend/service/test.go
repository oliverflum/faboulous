package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/internal/db"
	"github.com/oliverflum/faboulous/backend/model"
	"gorm.io/gorm"
)

// sendTestResponse handles the common logic for sending test responses
func SendTestResponse(c *fiber.Ctx, test *model.Test, statusCode int) error {
	payload, err := model.NewTestPayload(test)
	if err != nil {
		return err
	}
	return c.Status(statusCode).JSON(payload)
}

// getTestByID retrieves a test by ID and returns an error if not found
func GetTestByID(id uint, preloadVariants bool) (*model.Test, error) {
	var test model.Test
	var result *gorm.DB
	if preloadVariants {
		result = db.GetDB().
			Preload("Variants").
			Preload("Variants.Features").
			First(&test, id)
	} else {
		result = db.GetDB().First(&test, id)
	}
	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Test not found")
	}
	return &test, nil
}
