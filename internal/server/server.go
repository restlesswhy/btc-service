package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/restlesswhy/btc-service/config"
	v1 "github.com/restlesswhy/btc-service/internal/delivery/http/v1"
	"github.com/restlesswhy/btc-service/internal/ledger"
	"github.com/restlesswhy/btc-service/internal/service"
	"github.com/restlesswhy/btc-service/internal/synch"
	"github.com/restlesswhy/btc-service/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"

	"net/http"
	_ "net/http/pprof"
)

type server struct {
	log   logger.Logger
	cfg   *config.Config
	pool  *pgxpool.Pool
	fiber *fiber.App
}

func New(log logger.Logger, cfg *config.Config, pool *pgxpool.Pool) *server {
	return &server{log: log, cfg: cfg, pool: pool, fiber: fiber.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ledger := ledger.New(s.cfg, s.log)

	synch, err := synch.New(s.cfg, s.log, ledger)
	if err != nil {
		return errors.Wrap(err, "run synchronizer error")
	}
	defer synch.Close()

	service := service.New(s.log, ledger)
	controller := v1.New(s.log, service)
	controller.SetupRoutes(s.fiber)

	go func() {
		if err := s.runHttp(); err != nil {
			s.log.Errorf("(runHttp) err: %v", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: %v", s.getMicroserviceName(), s.cfg.Http.Port)

	go func() {
		s.log.Error(http.ListenAndServe(":6060", nil))
	}()

	<-ctx.Done()

	if err = s.fiber.Shutdown(); err != nil {
		s.log.Warnf("(Shutdown) err: %v", err)
		return err
	}

	s.log.Info("Service gracefully closed.")
	return nil
}

func (s server) getMicroserviceName() string {
	return fmt.Sprintf("(%s)", strings.ToUpper(s.cfg.ServiceName))
}
