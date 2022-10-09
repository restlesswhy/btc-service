package service

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/restlesswhy/btc-service/internal/ledger"
	"github.com/restlesswhy/btc-service/internal/models"
	"github.com/restlesswhy/btc-service/internal/utils"
	"github.com/restlesswhy/btc-service/pkg/logger"
)

type Ledger interface {
	NewInstance(id int) (ledger.Instance, error)
}

type service struct {
	log    logger.Logger
	ledger Ledger
}

func New(log logger.Logger, ledger Ledger) *service {
	return &service{log: log, ledger: ledger}
}

func (s *service) GetCurrency(symbol string) (models.Currency, error) {
	i, err := s.ledger.NewInstance(2)
	if err != nil {
		return models.Currency{}, errors.Wrap(err, "get new hf instance error")
	}
	defer i.Close()

	res, err := i.GetCurrency(strings.ToUpper(symbol))
	if err != nil {
		return models.Currency{}, errors.Wrap(err, "get currency error")
	}

	return res, nil
}

func (s *service) GetCurrencyPrice(symbol string) (models.QuotaDetailResponce, error) {
	i, err := s.ledger.NewInstance(2)
	if err != nil {
		return models.QuotaDetailResponce{}, errors.Wrap(err, "get new hf instance error")
	}
	defer i.Close()

	res, err := i.GetCurrencyPrice(strings.ToUpper(symbol))
	if err != nil {
		return models.QuotaDetailResponce{}, errors.Wrap(err, "get currency error")
	}

	return res.Convert(), nil
}

func (s *service) GetCurrencyPriceFromHistory(symbol string) ([]models.QuotaDetailResponce, error) {
	i, err := s.ledger.NewInstance(2)
	if err != nil {
		return nil, errors.Wrap(err, "get new hf instance error")
	}
	defer i.Close()

	quotes, err := i.GetCurrencyPriceFromHistory(strings.ToUpper(symbol))
	if err != nil {
		return nil, errors.Wrap(err, "get currency error")
	}

	res := make([]models.QuotaDetailResponce, 0, len(quotes))
	for _, v := range quotes {
		res = append(res, v.Convert())
	}	

	return res, nil
}

func (s *service) GetAllCurrentPrices() ([]models.QuotaDetailResponce, error) {
	i, err := s.ledger.NewInstance(2)
	if err != nil {
		return nil, errors.Wrap(err, "get new hf instance error")
	}
	defer i.Close()

	symbols, err := utils.GetQuotesSymbols()
	if err != nil {
		return nil, errors.Wrap(err, "get quotes symbols error")
	}

	quotes, err := i.GetAllCurrentPrices(symbols)
	if err != nil {
		return nil, errors.Wrap(err, "get all current prices error")
	}

	res := make([]models.QuotaDetailResponce, 0, len(quotes))
	for _, v := range quotes {
		res = append(res, v.Convert())
	}

	return res, nil
}
