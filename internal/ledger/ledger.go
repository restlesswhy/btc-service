package ledger

import (
	"sync"

	"github.com/restlesswhy/btc-service/config"
	"github.com/restlesswhy/btc-service/pkg/logger"
)

type instanceRequest struct {
	inst chan *instance
	err  chan error
	id   int
}

type ledger struct {
	wg  *sync.WaitGroup
	cfg *config.Config
	log logger.Logger

	newInstCh   chan instanceRequest
	closeInstCh chan int

	close chan struct{}
}

func New(cfg *config.Config, log logger.Logger) *ledger {
	l := &ledger{
		wg:          &sync.WaitGroup{},
		newInstCh:   make(chan instanceRequest),
		close:       make(chan struct{}),
		closeInstCh: make(chan int),
		cfg:         cfg,
		log:         log,
	}

	l.wg.Add(1)
	go l.run()

	return l
}

func (l *ledger) run() {
	defer l.wg.Done()

	pool := make(map[int]*instance)

ledger:
	for {
		select {
		case <-l.close:
			break ledger

		case req := <-l.newInstCh:
			l.log.Info("starting new unstance")
			if inst, ok := pool[req.id]; ok {
				req.err <- nil
				req.inst <- inst
				continue
			}

			i, err := newInstance(l.cfg, req.id, l.closeInstCh, l.log)
			if err != nil {
				req.err <- err
				continue
			}

			pool[req.id] = i

			req.err <- nil
			req.inst <- i
			l.log.Info("new instance added")

		case id := <-l.closeInstCh:
			delete(pool, id)
			l.log.Info("instance closed")
		}
	}
}

func (l *ledger) NewInstance(id int) (Instance, error) {
	req := instanceRequest{
		id:   id,
		inst: make(chan *instance),
		err:  make(chan error),
	}
	defer close(req.err)
	defer close(req.inst)

	l.newInstCh <- req

	err := <-req.err
	if err != nil {
		return nil, err
	}

	return <-req.inst, nil
}

func (l *ledger) Close() error {
	close(l.close)
	l.wg.Wait()

	return nil
}
