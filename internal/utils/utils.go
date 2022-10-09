package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/restlesswhy/btc-service/internal/models"
)

func LoadQoutes() ([]models.QuotaDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	client := http.DefaultClient

	req, _ := http.NewRequest(
		"GET", "https://www.blockchain.com/ru/ticker", nil,
	)

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "send request error")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read request body error")
	}

	q := make(models.Quotes, 0)
	if err := json.Unmarshal(body, &q); err != nil {
		return nil, errors.Wrap(err, "unmarshal request body error")
	}

	return q.Convert()
}

func GetQuotesSymbols() ([]string, error) {
	q, err := LoadQoutes()
	if err != nil {
		return nil, errors.Wrap(err, "load currencies error")
	}

	res := make([]string, 0, len(q))
	for _, v := range q {
		res = append(res, v.Symbol)
	}

	return res, nil
}

func LoadCurrencies() ([]models.Currency, error) {
	f, err := os.ReadFile(`docs/en_currencies.json`)
	if err != nil {
		return nil, errors.Wrap(err, "read en curr error")
	}

	en := make([]models.ENCurrencyInfo, 0)
	if err := json.Unmarshal(f, &en); err != nil {
		return nil, errors.Wrap(err, "unmarshal request body error")
	}

	f, err = os.ReadFile(`docs/ru_currencies.json`)
	if err != nil {
		return nil, errors.Wrap(err, "read ru curr error")
	}

	ru := make([]models.RUCurrencyInfo, 0)
	if err := json.Unmarshal(f, &ru); err != nil {
		return nil, errors.Wrap(err, "unmarshal request body error")
	}

	res := make([]models.Currency, 0, len(en))
	for _, e := range en {
		curr := models.Currency{
			Symbol: e.Symbol,
			ENName: e.Name,
		}

		for _, r := range ru {
			if strings.TrimSpace(e.Symbol) == strings.TrimSpace(r.Symbol) {
				curr.Code = r.Code
				curr.RUName = r.Name
				break
			}
		}

		if curr.Code == "" {
			continue
		}

		res = append(res, curr)
	}

	return res, nil
}
