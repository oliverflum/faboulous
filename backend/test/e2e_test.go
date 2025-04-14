package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oliverflum/faboulous/backend/app"
	"github.com/oliverflum/faboulous/backend/model"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	os.Setenv("FAB_DB_TYPE", "memory")
	os.Setenv("FAB_ENV", "test")

	app := app.SetupApp()

	return app
}

func TestFeatureEndpoints(t *testing.T) {
	app := setupTestApp()

	// Test Create Feature
	t.Run("Create Feature", func(t *testing.T) {
		feature := model.FeatureWritePayload{
			Name:  "Test Feature",
			Value: "Test Value",
		}
		body, _ := json.Marshal(feature)

		req := httptest.NewRequest("POST", "/api/config/feature/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	// Test List Features
	t.Run("List Features", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/config/feature/", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Get Feature
	t.Run("Get Feature", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/config/feature/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Update Feature
	t.Run("Update Feature", func(t *testing.T) {
		feature := model.FeatureWritePayload{
			Name:  "Updated Feature",
			Value: "Updated Value",
		}
		body, _ := json.Marshal(feature)

		req := httptest.NewRequest("PUT", "/api/config/feature/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Delete Feature
	t.Run("Delete Feature", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/config/feature/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestTestEndpoints(t *testing.T) {
	app := setupTestApp()

	// Test Create Test
	t.Run("Create Test", func(t *testing.T) {
		test := model.TestWritePayload{
			Name:   "Test AB Test",
			Active: true,
			Method: model.HASH,
		}
		body, _ := json.Marshal(test)

		req := httptest.NewRequest("POST", "/api/config/test/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	// Test List Tests
	t.Run("List Tests", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/config/test/", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Get Test
	t.Run("Get Test", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/config/test/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Update Test
	t.Run("Update Test", func(t *testing.T) {
		test := model.TestWritePayload{
			Name:   "Updated Test",
			Active: false,
			Method: model.RANDOM,
		}
		body, _ := json.Marshal(test)

		req := httptest.NewRequest("PUT", "/api/config/test/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Delete Test
	t.Run("Delete Test", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/config/test/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestVariantEndpoints(t *testing.T) {
	app := setupTestApp()

	// First create a test
	test := model.TestWritePayload{
		Name:   "Test AB Test",
		Active: true,
		Method: model.HASH,
	}
	testBody, _ := json.Marshal(test)
	testReq := httptest.NewRequest("POST", "/api/config/test/", bytes.NewBuffer(testBody))
	testReq.Header.Set("Content-Type", "application/json")
	testResp, err := app.Test(testReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, testResp.StatusCode)

	// Create two features
	feature1 := model.FeatureWritePayload{
		Name:  "Test Feature 1",
		Value: "Value 1",
	}
	feature1Body, _ := json.Marshal(feature1)
	feature1Req := httptest.NewRequest("POST", "/api/config/feature/", bytes.NewBuffer(feature1Body))
	feature1Req.Header.Set("Content-Type", "application/json")
	feature1Resp, err := app.Test(feature1Req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, feature1Resp.StatusCode)

	feature2 := model.FeatureWritePayload{
		Name:  "Test Feature 2",
		Value: 2,
	}
	feature2Body, _ := json.Marshal(feature2)
	feature2Req := httptest.NewRequest("POST", "/api/config/feature/", bytes.NewBuffer(feature2Body))
	feature2Req.Header.Set("Content-Type", "application/json")
	feature2Resp, err := app.Test(feature2Req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, feature2Resp.StatusCode)

	// Test Create Variant
	t.Run("Create Variant", func(t *testing.T) {
		variant := model.VariantWritePayload{
			Name: "Test Variant Premium",
		}
		body, _ := json.Marshal(variant)

		req := httptest.NewRequest("POST", "/api/config/test/1/variant/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	// Test Update Variant
	t.Run("Update Variant", func(t *testing.T) {
		variant := model.VariantWritePayload{
			Name: "Updated Variant",
		}
		body, _ := json.Marshal(variant)

		req := httptest.NewRequest("PUT", "/api/config/test/1/variant/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// Test Delete Variant
	t.Run("Delete Variant", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/config/test/1/variant/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}
