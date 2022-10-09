package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/restlesswhy/btc-service/internal/models"
	"github.com/restlesswhy/btc-service/pkg/logger"
)

type Service interface {
	GetCurrency(symbol string) (models.Currency, error)
	GetCurrencyPrice(symbol string) (models.QuotaDetailResponce, error)
	GetAllCurrentPrices() ([]models.QuotaDetailResponce, error)
	GetCurrencyPriceFromHistory(symbol string) ([]models.QuotaDetailResponce, error)
}

type handler struct {
	log     logger.Logger
	service Service
}

func New(log logger.Logger, service Service) *handler {
	return &handler{log: log, service: service}
}

// GetCurrency godoc
// @Summary Get currency info
// @Description send currency symbol, get info
// @Tags Currency
// @Accept json
// @Produce json
// @Param symbol query string false "Currency identificator"
// @Success 200 {object} models.Currency
// @Router /currency [get]
func (h *handler) getCurrency(c *fiber.Ctx) error {
	symbol := string(c.Context().FormValue("symbol"))
	if len(symbol) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "id parameter is empty",
		})
	}

	currency, err := h.service.GetCurrency(symbol)
	if err != nil {
		h.log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "get currency error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   currency,
	})
}

// GetCurrencyPrice godoc
// @Summary Get currency quota
// @Description send currency symbol, get currency quota
// @Tags Currency
// @Accept json
// @Produce json
// @Param symbol query string false "Currency identificator"
// @Success 200 {object} models.QuotaDetailResponce
// @Router /currency/price [get]
func (h *handler) getCurrencyPrice(c *fiber.Ctx) error {
	symbol := string(c.Context().FormValue("symbol"))
	if len(symbol) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "id parameter is empty",
		})
	}

	quota, err := h.service.GetCurrencyPrice(symbol)
	if err != nil {
		h.log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "get quota error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   quota,
	})
}

// GetCurrencyPriceFromHistory godoc
// @Summary Get currency quota history
// @Description send currency symbol, get info
// @Tags Currency
// @Accept json
// @Produce json
// @Param symbol query string false "Currency identificator"
// @Success 200 {object} []models.QuotaDetailResponce
// @Router /currency/price/history [get]
func (h *handler) getCurrencyPriceFromHistory(c *fiber.Ctx) error {
	symbol := string(c.Context().FormValue("symbol"))
	if len(symbol) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "id parameter is empty",
		})
	}

	quotes, err := h.service.GetCurrencyPriceFromHistory(symbol)
	if err != nil {
		h.log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "get quota from history error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   quotes,
	})
}

// GetAllCurrentPrices godoc
// @Summary Get all quotas of currencies
// @Description send currency symbol, get info
// @Tags Currency
// @Accept json
// @Produce json
// @Param symbol query string false "Currency identificator"
// @Success 200 {object} []models.QuotaDetailResponce
// @Router /currencies/price [get]
func (h *handler) getAllCurrentPrices(c *fiber.Ctx) error {
	quotes, err := h.service.GetAllCurrentPrices()
	if err != nil {
		h.log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "get all quotes error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   quotes,
	})
}
