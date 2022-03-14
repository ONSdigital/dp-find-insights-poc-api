package service

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/middleware"
	"github.com/ONSdigital/dp-find-insights-poc-api/api"
	"github.com/ONSdigital/dp-find-insights-poc-api/cache"
	"github.com/ONSdigital/dp-find-insights-poc-api/cantabular"
	"github.com/ONSdigital/dp-find-insights-poc-api/config"
	"github.com/ONSdigital/dp-find-insights-poc-api/handlers"
	"github.com/ONSdigital/dp-find-insights-poc-api/metadata"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/geodata"
	"github.com/ONSdigital/dp-find-insights-poc-api/postcode"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	var cant *cantabular.Client
	if cfg.EnableCantabular {
		cant = cantabular.New(cfg.CantabularURL, cfg.CantabularUser, os.Getenv("CANT_PW"))
	}

	var db *database.Database
	var queryGeodata *geodata.Geodata
	var md *metadata.Metadata
	var pc *postcode.Postcode
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

		db, err = database.Open("pgx", database.GetDSN(pgpwd))
		if err != nil {
			return nil, err
		}

		// set up our query functionality if we have a db
		queryGeodata, err = geodata.New(db, cant, cfg.MaxMetrics)
		if err != nil {
			return nil, err
		}

		// metadata.New can set up gorm itself, but it calls GetDSN without an
		// argument, so it cannot know about passwords held in AWS secrets.
		//
		// We loop here in case the db isn't up yet (happens when using docker compose).
		// (Looping doesn't seem to be needed for the pgx connection, for some reason.)
		var gdb *gorm.DB
		for try := 0; try < 5; try++ {
			gdb, err = gorm.Open(postgres.Open(database.GetDSN(pgpwd)), &gorm.Config{
				//	Logger: logger.Default.LogMode(logger.Info), // display SQL
			})
			if err == nil {
				break
			}
			log.Info(ctx, "cannot-connect-to-postgres-yet")
			time.Sleep(1 * time.Second)
		}
		if err != nil {
			return nil, err
		}

		md, err = metadata.New(gdb)
		if err != nil {
			return nil, err
		}

		pc = postcode.New(gdb)

	}

	cm, err := cache.New(cfg.CacheTTL, cfg.CacheSize)
	if err != nil {
		return nil, err
	}

	// Setup the API
	a := handlers.New(true, queryGeodata, md, cm, pc) // always include private handlers for now

	// Setup health checks
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}
	if err := registerCheckers(ctx, hc, db, md, cant); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}
	hc.Start(ctx)

	timeoutHandler := func(h http.Handler) http.Handler {
		return http.TimeoutHandler(h, cfg.WriteTimeout, "operation timed out\n")
	}

	// build handler chain
	chain := alice.New(
		middleware.Whitelist(middleware.HealthcheckFilter(hc.Handler)),
		timeoutHandler,
	).Then(api.Handler(a))

	// bind router handler to http server
	s := serviceList.GetHTTPServer(cfg.BindAddr, chain)

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
	md *metadata.Metadata,
	cant *cantabular.Client,
) (err error) {
	if db != nil {
		err = hc.AddCheck("postgres", db.Checker)
	}
	if md != nil {
		err = hc.AddCheck("gorm", md.Checker)
	}
	if cant != nil {
		err = hc.AddCheck("cantabular", cant.Checker)
	}
	return err
}
