package service

import (
	"context"
	"fmt"
	"os"

	"github.com/ONSdigital/dp-api-clients-go/middleware"
	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/config"
	"github.com/ONSdigital/dp-find-insights-poc-api/handlers"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/demo"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {

	log.Info(ctx, "running service")

	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	var db *database.Database
	var queryDemo *demo.Demo
	var err error
	if cfg.EnableDatabase {
		// figure out postgres password
		pgpwd := os.Getenv("PGPASSWORD")
		if pgpwd == "" {
			aws, err := serviceList.GetAWS()
			if err != nil {
				return nil, err
			}
			pgpwd, err = aws.GetSecret(os.Getenv("FI_PG_SECRET_ID"))
			if err != nil {
				return nil, err
			}
		}

		// open postgres connection
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			os.Getenv("PGUSER"),
			pgpwd,
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT"),
			os.Getenv("PGDATABASE"),
		)
		db, err = database.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}

		// set up our query functionality if we have a db
		queryDemo, err = demo.New(db)
		if err != nil {
			return nil, err
		}
	}

	// Setup the API
	a := handlers.New(true, queryDemo) // always include private handlers for now

	// Setup health checks
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}
	if err := registerCheckers(ctx, hc, db); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}
	hc.Start(ctx)

	// put health checker at beginning of middleware chain
	chain := alice.New(middleware.Whitelist(middleware.HealthcheckFilter(hc.Handler)))

	// attach the appropriate api to the chain to create full router
	rtr := chain.Then(api.Handler(a))

	// bind router handler to http server
	s := serviceList.GetHTTPServer(cfg.BindAddr, rtr)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      s,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	hc HealthChecker,
	db *database.Database,
) (err error) {

	// TODO: add other health checks here, as per dp-upload-service

	if db != nil {
		err = hc.AddCheck("postgres", db.Checker)
	}
	return nil
}
