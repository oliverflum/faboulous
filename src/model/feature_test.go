package model

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/util"
	"github.com/stretchr/testify/assert"
)

func TestGetEntityValueAndType(t *testing.T) {
	tests := []struct {
		name          string
		value         any
		expectedType  string
		expectedValue string
		expectedError error
	}{
		{
			name:          "String value",
			value:         "test",
			expectedType:  STRING,
			expectedValue: "test",
			expectedError: nil,
		},
		{
			name:          "Int value",
			value:         42,
			expectedType:  INT,
			expectedValue: "42",
			expectedError: nil,
		},
		{
			name:          "Bool value",
			value:         true,
			expectedType:  BOOL,
			expectedValue: "true",
			expectedError: nil,
		},
		{
			name:          "Float value",
			value:         3.14,
			expectedType:  FLOAT,
			expectedValue: "3.140000",
			expectedError: nil,
		},
		{
			name:          "Unsupported type",
			value:         struct{}{},
			expectedType:  "",
			expectedValue: "",
			expectedError: util.FabolousError{
				Code:    fiber.StatusBadRequest,
				Message: "unsupported value type: struct {}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valueType, stringValue, err := GetEntityValueAndType(tt.value)

			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedType, valueType)
				assert.Equal(t, tt.expectedValue, stringValue)
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

func TestGetPayloadValue(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		valueType     string
		expectedValue any
		expectedError error
	}{
		{
			name:          "String value",
			value:         "test",
			valueType:     STRING,
			expectedValue: "test",
			expectedError: nil,
		},
		{
			name:          "Int value",
			value:         "42",
			valueType:     INT,
			expectedValue: 42,
			expectedError: nil,
		},
		{
			name:          "Bool true value",
			value:         "true",
			valueType:     BOOL,
			expectedValue: true,
			expectedError: nil,
		},
		{
			name:          "Bool false value",
			value:         "false",
			valueType:     BOOL,
			expectedValue: false,
			expectedError: nil,
		},
		{
			name:          "Float value",
			value:         "3.14",
			valueType:     FLOAT,
			expectedValue: 3.14,
			expectedError: nil,
		},
		{
			name:          "Invalid int value",
			value:         "not_an_int",
			valueType:     INT,
			expectedValue: nil,
			expectedError: util.FabolousError{
				Code:    fiber.StatusInternalServerError,
				Message: "Invalid value for " + INT + ": " + "not_an_int",
			},
		},
		{
			name:          "Invalid float value",
			value:         "not_a_float",
			valueType:     FLOAT,
			expectedValue: nil,
			expectedError: util.FabolousError{
				Code:    fiber.StatusInternalServerError,
				Message: "Invalid value for " + FLOAT + ": " + "not_a_float",
			},
		},
		{
			name:          "Unsupported type",
			value:         "value",
			valueType:     "UNKNOWN",
			expectedValue: nil,
			expectedError: util.FabolousError{
				Code:    fiber.StatusInternalServerError,
				Message: "unsupported value type: UNKNOWN",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := GetPayloadValue(tt.value, tt.valueType)

			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
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
