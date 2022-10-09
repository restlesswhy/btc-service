package models

import (
	"time"

	"github.com/pkg/errors"
)

type Quotes map[string]map[string]any

type QuotaDetail struct {
	Buy    int64  `json:"buy"`
	Sell   int64  `json:"sell"`
	Symbol string `json:"symbol"`
	Time   int64  `json:"time"`
}

func (q *QuotaDetail) Convert() QuotaDetailResponce {
	return QuotaDetailResponce{
		Buy:    float64(q.Buy) / 100,
		Sell:   float64(q.Sell) / 100,
		Symbol: q.Symbol,
		Time:   time.Unix(q.Time, 0).Format(time.RFC3339),
	}
}

type QuotaDetailResponce struct {
	Buy    float64 `json:"buy"`
	Sell   float64 `json:"sell"`
	Symbol string  `json:"symbol"`
	Time   string  `json:"time"`
}

func (q Quotes) Convert() ([]QuotaDetail, error) {
	res := make([]QuotaDetail, 0, len(q))

	for _, v := range q {
		quota := QuotaDetail{}
		for k, data := range v {
			switch k {
			case "buy":
				s, ok := data.(float64)
				if !ok {
					return nil, errors.New("convert error, 'buy' field is not float64")
				}

				quota.Buy = int64(s * 100)

			case "sell":
				s, ok := data.(float64)
				if !ok {
					return nil, errors.New("convert error, 'sell' field is not float64")
				}

				quota.Sell = int64(s * 100)

			case "symbol":
				s, ok := data.(string)
				if !ok {
					return nil, errors.New("convert error, 'symbol' field is not string")
				}

				quota.Symbol = s
			}
		}

		quota.Time = time.Now().Unix()
		res = append(res, quota)
	}

	return res, nil
}

type Currency struct {
	Symbol string `json:"symbol"`
	RUName string `json:"ru_name"`
	ENName string `json:"en_name"`
	Code   string `json:"code"`
}

type ENCurrencyInfo struct {
	Symbol string `json:"cc"`
	Name   string `json:"name"`
}

type RUCurrencyInfo struct {
	Symbol string `json:"STRCODE"`
	Name   string `json:"NAME"`
	Code   string `json:"CODE"`
}
