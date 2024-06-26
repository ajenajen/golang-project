//go:build integration

package handlers_test

import (
	"fmt"
	"gotest/handlers"
	"gotest/repositories"
	"gotest/services"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPromotionCalculateDiscountIntegrationService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		amount := 100
		expected := 80

		promoRepo := repositories.NewPromotionRepositoryMock()
		promoRepo.On("GetPromotion").Return(repositories.Promotion{
			ID:              1,
			PurchaseMin:     100,
			DiscountPercent: 20,
		}, nil)

		promoService := services.NewPromotionService(promoRepo)
		promoHandler := handlers.NewPromotionHandler(promoService)

		// http://localhost:8000/calculate?amount=100
		app := fiber.New()
		app.Get("/calculate", promoHandler.CalculateDiscount)

		req := httptest.NewRequest("GET", fmt.Sprintf("/calculate?amount=%v", amount), nil)

		// Act
		res, _ := app.Test(req)
		defer res.Body.Close()

		// Assert
		if assert.Equal(t, fiber.StatusOK, res.StatusCode) { // ถ้าอ่านแล้วไม่ equal มันจะไม่มี body เลยต้องดักไว้
			body, _ := io.ReadAll(res.Body) //res.Body เป็น io.ReadCloser, body return เป็น slice of byte
			assert.Equal(t, strconv.Itoa(expected), string(body))
		}
	})
}
