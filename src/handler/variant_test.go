package handler

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/model"
	"github.com/oliverflum/faboulous/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.Feature{}, &model.Variant{}, &model.VariantFeature{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUpdateVariantFeatureValue(t *testing.T) {
	db := setupTestDB(t)

	// Create test features with different types
	features := []model.Feature{
		{
			Name:         "string_feature",
			Type:         model.STRING,
			DefaultValue: "default",
		},
		{
			Name:         "int_feature",
			Type:         model.INT,
			DefaultValue: "0",
		},
		{
			Name:         "bool_feature",
			Type:         model.BOOL,
			DefaultValue: "false",
		},
		{
			Name:         "float_feature",
			Type:         model.FLOAT,
			DefaultValue: "0.0",
		},
	}

	for _, feature := range features {
		err := db.Create(&feature).Error
		assert.NoError(t, err)
	}

	variant := model.Variant{
		Name: "test_variant",
	}
	err := db.Create(&variant).Error
	assert.NoError(t, err)

	tests := []struct {
		name          string
		feature       model.FeaturePayload
		value         any
		expectedError error
	}{
		{
			name: "Valid string value",
			feature: model.FeaturePayload{
				Name:  "string_feature",
				Value: "new_value",
			},
			value:         "new_value",
			expectedError: nil,
		},
		{
			name: "Valid int value",
			feature: model.FeaturePayload{
				Name:  "int_feature",
				Value: 42,
			},
			value:         42,
			expectedError: nil,
		},
		{
			name: "Valid bool value",
			feature: model.FeaturePayload{
				Name:  "bool_feature",
				Value: true,
			},
			value:         true,
			expectedError: nil,
		},
		{
			name: "Valid float value",
			feature: model.FeaturePayload{
				Name:  "float_feature",
				Value: 3.14,
			},
			value:         3.14,
			expectedError: nil,
		},
		{
			name: "Type mismatch - string to int",
			feature: model.FeaturePayload{
				Name:  "int_feature",
				Value: "not_an_int",
			},
			value: "not_an_int",
			expectedError: util.FabolousError{
				Code:    fiber.StatusBadRequest,
				Message: "value type mismatch: INT != STRING",
			},
		},
		{
			name: "Non-existent feature",
			feature: model.FeaturePayload{
				Name:  "non_existent",
				Value: "value",
			},
			value: "value",
			expectedError: util.FabolousError{
				Code:    fiber.StatusNotFound,
				Message: "Record not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateVariantFeatureValue(db, &variant, &tt.feature, tt.value)

			if tt.expectedError == nil {
				assert.NoError(t, err)

				// Find the feature by name
				var feature model.Feature
				err = db.Where("name = ?", tt.feature.Name).First(&feature).Error
				assert.NoError(t, err)

				// Verify the value was saved correctly
				var variantFeature model.VariantFeature
				err = db.Where("variant_id = ? AND feature_id = ?", variant.ID, feature.ID).First(&variantFeature).Error
				assert.NoError(t, err)

				// Convert the stored string value back to the expected type
				actualValue, err := model.GetPayloadValue(variantFeature.Value, feature.Type)
				assert.NoError(t, err)
				assert.Equal(t, tt.value, actualValue)
			} else {
				assert.Error(t, err)
				if fabError, ok := tt.expectedError.(util.FabolousError); ok {
					assert.Equal(t, fabError.Code, err.(util.FabolousError).Code)
					assert.Equal(t, fabError.Message, err.(util.FabolousError).Message)
				} else {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			}
		})
	}
}
