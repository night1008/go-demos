package main

import (
	"context"
	"fmt"
	"gin-api-server/internal/api"
	"gin-api-server/internal/app"
	"gin-api-server/internal/authority"
	"gin-api-server/internal/config"
	"gin-api-server/internal/router"
	"gin-api-server/internal/service"

	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

// @title FunnyDB
// @version 0.0.1
// @description funnydb api doc
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http
// @basePath /api/v1
func main() {
	app := &cli.App{
		Name:  "funnydb website",
		Usage: "funnydb website cli entrance",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Value:   "config.yml",
				Aliases: []string{"c"},
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.Load(c.String("config"))
			if err != nil {
				return fmt.Errorf("Fatal error config file: %w \n", err)
			}

			eg, ctx := errgroup.WithContext(context.Background())

			logger := app.Logger{Cfg: &cfg.Logger}
			if err := logger.Setup(); err != nil {
				return err
			}

			dbApp := app.DB{Cfg: &cfg.Database}
			db, err := dbApp.Setup()
			if err != nil {
				return err
			}
			seeder := app.Seeder{Cfg: &cfg.Seeder, DB: db}
			if err := seeder.Setup(); err != nil {
				return err
			}

			auther := &app.PwdAuther{
				DB: db,
			}
			router := &router.Router{
				Auther: auther,
				LoginAPI: &api.LoginAPI{
					UserSvc: &service.UserSvc{
						DB: db,
					},
				},
				RegisterAPI: &api.RegisterAPI{
					UserSvc: &service.UserSvc{
						DB: db,
					},
				},
				UserAPI: &api.UserAPI{
					UserSvc: &service.UserSvc{
						DB:        db,
						Authority: &authority.Authority{DB: db},
					},
				},
				OrgAPI: &api.OrgAPI{
					OrgSvc: &service.OrgSvc{
						DB:        db,
						Authority: &authority.Authority{DB: db},
					},
				},
				AppAPI: &api.AppAPI{
					AppSvc: &service.AppSvc{
						DB:        db,
						Authority: &authority.Authority{DB: db},
					},
				},
				RoleAPI: &api.RoleAPI{
					RoleSvc: &service.RoleSvc{
						DB: db,
					},
				},
			}

			httpServer := app.HTTPServer{
				Cfg:    &cfg.HTTPServer,
				Router: router,
			}
			if err := httpServer.Setup(); err != nil {
				return err
			}
			if err := httpServer.Run(ctx, eg); err != nil {
				return err
			}

			eg.Go(func() error {
				return terminatedSignalReceived(ctx, os.Interrupt, syscall.SIGTERM)
			})

			if err := eg.Wait(); err != nil {
				log.Info().Err(err).Msg("")
			}

			return nil
		},
	}
	app.RunAndExitOnError()
}

func terminatedSignalReceived(ctx context.Context, sigs ...os.Signal) error {
	quit := make(chan os.Signal, 0)
	signal.Notify(quit, sigs...)
	defer signal.Reset(sigs...)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-quit:
		return fmt.Errorf("received %v signal", sig)
	}
}
