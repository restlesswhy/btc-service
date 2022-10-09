package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/restlesswhy/btc-service/swag"
	swagger "github.com/arsmn/fiber-swagger/v2"
)

func (h *handler) SetupRoutes(r *fiber.App) {
	api := r.Group("/api/v1", logger.New())

	api.Get("/currency", h.getCurrency)
	api.Get("/currency/price", h.getCurrencyPrice)
	api.Get("/currency/price/history", h.getCurrencyPriceFromHistory)
	api.Get("/currencies/price", h.getAllCurrentPrices)
	api.Get("/swagger/*", swagger.HandlerDefault)
}
