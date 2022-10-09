package ledger

import (
	"context"
	"encoding/json"
	"io"

	"github.com/restlesswhy/btc-service/internal/models"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	gwproto "github.com/hyperledger/fabric-protos-go/gateway"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

type Instance interface {
	PutCurrency(data any) error
	GetCurrency(symbol string) (models.Currency, error)
	PutCurrencyPrice(data any) error
	GetCurrencyPrice(symbol string) (models.QuotaDetail, error)
	GetAllCurrentPrices(data any) ([]models.QuotaDetail, error)
	GetCurrencyPriceFromHistory(symbol string) ([]models.QuotaDetail, error)
	io.Closer
}

func (i *instance) PutCurrency(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal data error")
	}

	_, err = i.contract.SubmitTransaction("PutCurrency", string(b))
	if err != nil {
		return i.errorHandling(err)
	}

	return nil
}

func (i *instance) GetCurrency(symbol string) (models.Currency, error) {
	b, err := i.contract.EvaluateTransaction("GetCurrency", symbol)
	if err != nil {
		return models.Currency{}, i.errorHandling(err)
	}

	res := models.Currency{}
	if err := json.Unmarshal(b, &res); err != nil {
		return models.Currency{}, errors.Wrap(err, "unmarshal body error")
	}

	return res, nil
}

func (i *instance) PutCurrencyPrice(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal data error")
	}

	_, err = i.contract.SubmitTransaction("PutCurrencyPrice", string(b))
	if err != nil {
		return i.errorHandling(err)
	}

	return nil
}

func (i *instance) GetCurrencyPrice(symbol string) (models.QuotaDetail, error) {
	b, err := i.contract.EvaluateTransaction("GetCurrencyPrice", symbol)
	if err != nil {
		return models.QuotaDetail{}, i.errorHandling(err)
	}

	res := models.QuotaDetail{}
	if err := json.Unmarshal(b, &res); err != nil {
		return models.QuotaDetail{}, errors.Wrap(err, "unmarshal body error")
	}

	return res, nil
}

func (i *instance) GetAllCurrentPrices(data any) ([]models.QuotaDetail, error) {
	currencies, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "marshal currencies error")
	}
	
	b, err := i.contract.EvaluateTransaction("GetAllCurrentPrices", string(currencies))
	if err != nil {
		return nil, i.errorHandling(err)
	}

	quotes := []models.QuotaDetail{}
	if err := json.Unmarshal(b, &quotes); err != nil {
		return nil, errors.Wrap(err, "unmarshal body error")
	}

	return quotes, nil
}

func (i *instance) GetCurrencyPriceFromHistory(symbol string) ([]models.QuotaDetail, error) {
	b, err := i.contract.EvaluateTransaction("GetCurrencyPriceFromHistory", symbol)
	if err != nil {
		return nil, i.errorHandling(err)
	}

	quotas := []models.QuotaDetail{}
	if err := json.Unmarshal(b, &quotas); err != nil {
		return nil, errors.Wrap(err, "unmarshal body error")
	}

	return quotas, nil
}

func (i *instance) errorHandling(err error) error {
	switch err := err.(type) {
	case *client.EndorseError:
		i.log.Errorf("endorse error with gRPC status %v: %s", status.Code(err), err)
	case *client.SubmitError:
		i.log.Errorf("submit error with gRPC status %v: %s", status.Code(err), err)
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			i.log.Errorf("timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
		} else {
			i.log.Errorf("error obtaining commit status with gRPC status %v: %s", status.Code(err), err)
		}
	case *client.CommitError:
		i.log.Errorf("transaction %s failed to commit with status %d: %s", err.TransactionID, int32(err.Code), err)
	}

	statusErr := status.Convert(err)
	for _, detail := range statusErr.Details() {
		errDetail := detail.(*gwproto.ErrorDetail)
		i.log.Errorf("error from endpoint: %s, mspId: %s, message: %s", errDetail.Address, errDetail.MspId, errDetail.Message)
	}

	return err
}
