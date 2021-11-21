package app

import (
	"context"
	"crypto/tls"
	"gin-api-server/internal/config"
	"gin-api-server/internal/router"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type HTTPServer struct {
	Router *router.Router
	Cfg    *config.HTTPServerCfg
	Server *http.Server
}

func (hs *HTTPServer) Setup() error {
	hs.Router.Setup(hs.Cfg)

	mux := http.NewServeMux()
	mux.Handle("/api/", hs.Router.Engine)

	hs.Server = &http.Server{
		Addr:         hs.Cfg.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return nil
}

func (hs *HTTPServer) Run(ctx context.Context, eg *errgroup.Group) error {
	eg.Go(func() error {
		log.Info().Str("addr", hs.Server.Addr).Msg("listen on")
		var err error
		if hs.Cfg.CertFile != "" && hs.Cfg.KeyFile != "" {
			hs.Server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = hs.Server.ListenAndServeTLS(hs.Cfg.CertFile, hs.Cfg.KeyFile)
		} else {
			err = hs.Server.ListenAndServe()
		}
		return err
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			log.Info().Msg("start shutdown server...")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(hs.Cfg.ShutdownTimeout)*time.Second)
		defer cancel()
		hs.Server.SetKeepAlivesEnabled(false)
		var err error
		if err = hs.Server.Shutdown(shutdownCtx); err != nil {
			log.Info().Err(err).Msg("error during shutdown server")
		}
		return err
	})

	return nil
}
