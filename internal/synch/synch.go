package synch

import (
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/restlesswhy/btc-service/config"
	"github.com/restlesswhy/btc-service/internal/ledger"
	"github.com/restlesswhy/btc-service/internal/models"
	"github.com/restlesswhy/btc-service/internal/utils"
	"github.com/restlesswhy/btc-service/pkg/logger"
)

type Ledger interface {
	NewInstance(id int) (ledger.Instance, error)
}

type Synch interface {
	io.Closer
}

type synch struct {
	wg     *sync.WaitGroup
	cfg    *config.Config
	log    logger.Logger
	ledger Ledger

	close chan struct{}
}

func New(cfg *config.Config, log logger.Logger, ledger Ledger) (Synch, error) {
	synch := &synch{
		wg:     &sync.WaitGroup{},
		close:  make(chan struct{}),
		cfg:    cfg,
		log:    log,
		ledger: ledger,
	}

	if err := synch.storeCurrencies(); err != nil {
		return nil, errors.Wrap(err, "store currencies error")
	}

	synch.wg.Add(1)
	go synch.run()

	return synch, nil
}

func (s *synch) run() {
	defer s.wg.Done()

	t := time.NewTicker(s.cfg.Synch.Interval * time.Millisecond)

synch:
	for {
		select {
		case <-s.close:
			break synch

		case <-t.C:
			q, err := utils.LoadQoutes()
			if err != nil {
				s.log.Errorf("failed load qoutes: %v", err)
				continue
			}

			i, err := s.ledger.NewInstance(1)
			if err != nil {
				s.log.Errorf("create new ledger instance error: %v", err)
				continue
			}

			wg := sync.WaitGroup{}
			wg.Add(len(q))
			for _, v := range q {
				go func(quota models.QuotaDetail) {
					defer wg.Done()
					if err := i.PutCurrencyPrice(quota); err != nil {
						s.log.Errorf("put currency price error: %v", err)
					}
				}(v)
			}
			wg.Wait()

			i.Close()
		}
	}
}

func (s *synch) storeCurrencies() error {
	s.log.Debug("starting store currencies")

	curr, err := utils.LoadCurrencies()
	if err != nil {
		return errors.Wrap(err, "load currencies error")
	}

	i, err := s.ledger.NewInstance(1)
	if err != nil {
		return errors.Wrap(err, "create new ledger instance error")
	}
	defer i.Close()

	wg := sync.WaitGroup{}
	for _, v := range curr {
		wg.Add(1)
		go func(d models.Currency) {
			defer wg.Done()
			if err := i.PutCurrency(d); err != nil {
				s.log.Errorf("store currency error: %v", err)
			}

			s.log.Debugf("currency '%v' stored", d)
		}(v)
	}
	wg.Wait()

	s.log.Debug("stop store curr")
	return nil
}

func (s *synch) Close() error {
	close(s.close)
	s.wg.Wait()

	return nil
}
