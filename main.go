/*
This service provides authentication and authorization via OAuth2.

Copyright (C) 2020 Lars Gaubisch

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/go-utils/httputils"
	"github.com/rebel-l/smis"
	"github.com/rebel-l/smis/middleware/cors"

	"github.com/rebel-l/auth-service/auth"
	"github.com/rebel-l/auth-service/bootstrap"
	"github.com/rebel-l/auth-service/config"
	"github.com/rebel-l/auth-service/endpoint/doc"
	"github.com/rebel-l/auth-service/endpoint/facebook"
	"github.com/rebel-l/auth-service/endpoint/ping"
	"github.com/rebel-l/auth-service/endpoint/user"

	"github.com/sirupsen/logrus"
)

const (
	defaultPort      = 3000
	defaultTimeout   = 15
	defaultRedisAddr = "redis:6379"
	version          = "v0.1.0"
)

var (
	// nolint: godox
	closers      []io.Closer // TODO: add to code generator
	db           *sqlx.DB
	kv           *redis.Client
	log          logrus.FieldLogger
	port         *int
	tokenManager *auth.Manager
	svc          *smis.Service
)

func initCustomFlags() {
	/**
	  1. Add your custom service flags below, for more details see https://golang.org/pkg/flag/
	*/
}

func initCustom() error {
	/**
	  2. add your custom service initialisation below, e.g. database connection, caches etc.
	*/

	// Redis
	kv = redis.NewClient(&redis.Options{
		// nolint:godox
		Addr: defaultRedisAddr, // TODO: take from config
	})

	_, err := kv.Ping().Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	closers = append(closers, kv)

	// token Manager
	// nolint:godox
	tokenManager = auth.NewManager("secret1", "secret2", kv) // TODO: secrets must be injected by config / env

	// Middleware
	c := cors.Config{
		AccessControlAllowHeaders: []string{"*"},
		AccessControlAllowOrigins: []string{"https://www.shopfriend.test"},
		AccessControlMaxAge:       cors.AccessControlMaxAgeDefault,
	}
	// nolint:godox
	// TODO: init based on config, config should be loaded from specific file

	svc.WithDefaultMiddleware(c)

	authMiddleware, err := auth.NewMiddleware(svc, tokenManager)
	if err != nil {
		return fmt.Errorf("failed to setup auth middleware: %w", err)
	}
	svc.AddMiddlewareForRestrictedChain(authMiddleware.Handler)

	// Database
	// nolint:godox
	db, err = bootstrap.Database(&config.Database{}, version, true) // TODO: take connection from config
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	closers = append(closers, db)

	return nil
}

func initCustomRoutes() error {
	/**
	  3. Register your custom routes below
	*/
	if err := facebook.Init(svc, db, tokenManager, httputils.NewClient()); err != nil {
		return fmt.Errorf("failed to init facebook endpoint: %w", err)
	}

	if err := user.Init(svc, tokenManager); err != nil {
		return fmt.Errorf("failed to init user endpoint: %w", err)
	}

	return nil
}

func main() {
	log = logrus.New()
	log.Info("Starting service: auth-service")

	defer func() {
		for _, v := range closers {
			_ = v.Close()
		}
	}()

	initFlags()
	initService()

	if err := initCustom(); err != nil {
		log.Errorf("Failed to initialise custom settings: %s", err)

		return
	}

	if err := initRoutes(); err != nil {
		log.Errorf("Failed to initialise routes: %s", err)

		return
	}

	log.Infof("Service listens to port %d", *port)
	if err := svc.ListenAndServe(); err != nil {
		log.Errorf("Failed to start server: %s", err)

		return
	}
}

func initService() {
	router := mux.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", *port),
		WriteTimeout: defaultTimeout * time.Second,
		ReadTimeout:  defaultTimeout * time.Second,
	}

	var err error
	svc, err = smis.NewService(srv, router, log)
	if err != nil {
		log.Fatalf("failed to initialize service: %s", err)
	}
}

func initRoutes() error {
	if err := initDefaultRoutes(); err != nil {
		return fmt.Errorf("default routes failed: %w", err)
	}

	if err := initCustomRoutes(); err != nil {
		return fmt.Errorf("custom routes failed: %w", err)
	}

	return nil
}

func initDefaultRoutes() error {
	if err := ping.Init(svc); err != nil {
		return fmt.Errorf("failed to init endpoint /ping: %w", err)
	}

	if err := doc.Init(svc); err != nil {
		return fmt.Errorf("failed to init endpoint /doc: %w", err)
	}

	return nil
}

func initFlags() {
	initDefaultFlags()
	initCustomFlags()
	flag.Parse()
}

func initDefaultFlags() {
	port = flag.Int("p", defaultPort, "the port the service listens to")
}
