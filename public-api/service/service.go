package service

import (
	"context"
	"github.com/cafo13/animal-facts/pkg/router"
	"github.com/pkg/errors"
	"net/http"

	"github.com/neko-neko/echo-logrus/v2/log"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	router *router.Router
}

func NewService(router *router.Router) *Service {
	return &Service{router}
}

func (s *Service) Run(ctx context.Context, port int) error {
	log.Logger().Info("Server starting...")

	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		err := s.router.Run(port)

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.router.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
