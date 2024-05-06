package services_test

import (
	"errors"
	"gotest/repositories"
	"gotest/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromotionCalculateDiscount(t *testing.T) {

	type testCase struct {
		name            string
		purchaseMin     int
		discountPercent int
		amount          int
		expected        int
	}

	cases := []testCase{
		{name: "applied amount 100", purchaseMin: 100, discountPercent: 20, amount: 100, expected: 80},
		{name: "applied amount 200", purchaseMin: 100, discountPercent: 20, amount: 200, expected: 160},
		{name: "applied amount 300", purchaseMin: 100, discountPercent: 20, amount: 300, expected: 240},
		{name: "not applied amount 50", purchaseMin: 100, discountPercent: 20, amount: 50, expected: 50},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Arrange
			promoRepo := repositories.NewPromotionRepositoryMock()
			promoRepo.On("GetPromotion").Return(repositories.Promotion{
				ID:              1,
				PurchaseMin:     c.purchaseMin,
				DiscountPercent: c.discountPercent,
			}, nil)

			promoService := services.NewPromotionService(promoRepo)

			// Act
			actual, _ := promoService.CalculateDiscount(c.amount)
			expected := c.expected

			// Assert
			assert.Equal(t, expected, actual)
		})
	}

	t.Run("purchase zero amount", func(t *testing.T) {
		// Arrange
		promoRepo := repositories.NewPromotionRepositoryMock()
		promoRepo.On("GetPromotion").Return(repositories.Promotion{
			ID:              1,
			PurchaseMin:     100,
			DiscountPercent: 20,
		}, nil)

		promoService := services.NewPromotionService(promoRepo)

		// Act
		_, actualErr := promoService.CalculateDiscount(0)

		// Assert
		assert.ErrorIs(t, services.ErrZeroAmount, actualErr)
		promoRepo.AssertNotCalled(t, "GetPromotion") //กันเรื่องลำดับการเรียก ไม่ให้เรียก GetPromotion
	})

	t.Run("repository error", func(t *testing.T) {
		// Arrange
		promoRepo := repositories.NewPromotionRepositoryMock()
		promoRepo.On("GetPromotion").Return(repositories.Promotion{}, errors.New(""))

		promoService := services.NewPromotionService(promoRepo)

		// Act
		_, actualErr := promoService.CalculateDiscount(100)

		// Assert
		assert.ErrorIs(t, services.ErrRepository, actualErr)
	})
}
