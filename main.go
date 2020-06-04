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
	"net/http"
	"path"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rebel-l/auth-service/endpoint/doc"
	"github.com/rebel-l/auth-service/endpoint/facebook"
	"github.com/rebel-l/auth-service/endpoint/ping"
	"github.com/rebel-l/go-utils/httputils"
	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/schema"
	"github.com/rebel-l/smis"
	"github.com/rebel-l/smis/middleware/cors"

	"github.com/sirupsen/logrus"
)

const (
	defaultPort    = 3000
	defaultTimeout = 15
	sqlScriptPath  = "./scripts/sql"
	sqlStoragePath = "./storage"
	sqlStorageFile = sqlStoragePath + "/auth-service.db"
	version        = "v0.1.0"
)

var (
	db   *sqlx.DB
	log  logrus.FieldLogger
	port *int
	svc  *smis.Service
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

	// Middleware
	c := cors.Config{
		AccessControlAllowHeaders: []string{"*"},
		AccessControlAllowOrigins: []string{"https://www.shopfriend.test"},
		AccessControlMaxAge:       cors.AccessControlMaxAgeDefault,
	}
	// nolint:godox
	//TODO: init based on config, config should be loaded from specific file

	svc.WithDefaultMiddleware(c)

	// Database
	var err error
	if err = osutils.CreateDirectoryIfNotExists(sqlStoragePath); err != nil {
		return err
	}

	if err = osutils.CreateFileIfNotExists(sqlStorageFile); err != nil {
		return err
	}

	// nolint:godox
	db, err = sqlx.Open("sqlite3", sqlStorageFile) // TODO: take connection from config
	if err != nil {
		return err
	}

	s := schema.New(db)
	s.WithProgressBar()
	if err = s.Upgrade(path.Join(sqlScriptPath, "sqlite"), version); err != nil {
		return err
	}

	return nil
}

func initCustomRoutes() error {
	/**
	  3. Register your custom routes below
	*/
	if err := facebook.Init(svc, db, httputils.NewClient()); err != nil {
		return fmt.Errorf("failed to init facebook endpoint: %w", err)
	}

	return nil
}

func main() {
	log = logrus.New()
	log.Info("Starting service: auth-service")

	initFlags()
	initService()

	if err := initCustom(); err != nil {
		log.Fatalf("Failed to initialise custom settings: %s", err)
	}

	if err := initRoutes(); err != nil {
		log.Fatalf("Failed to initialise routes: %s", err)
	}

	log.Infof("Service listens to port %d", *port)
	if err := svc.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
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
		return err
	}

	if err := doc.Init(svc); err != nil {
		return err
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
