// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package exportimport implements the handlers for the export-importer functionality.
package exportimport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/exposure-notifications-server/internal/database"
	eidb "github.com/google/exposure-notifications-server/internal/exportimport/database"
	pubdb "github.com/google/exposure-notifications-server/internal/publish/database"
	"github.com/google/exposure-notifications-server/internal/serverenv"
	"github.com/google/exposure-notifications-server/pkg/server"
)

// Server hosts end points to manage key rotation
type Server struct {
	config         *Config
	env            *serverenv.ServerEnv
	db             *database.DB
	exportImportDB *eidb.ExportImportDB
	publishDB      *pubdb.PublishDB
}

// NewServer creates a Server that manages deletion of
// old export files that are no longer needed by clients for download.
func NewServer(config *Config, env *serverenv.ServerEnv) (*Server, error) {
	if env.Database() == nil {
		return nil, fmt.Errorf("missing database in server environment")
	}

	db := env.Database()
	exportImportDB := eidb.New(db)
	publishDB := pubdb.New(db)

	return &Server{
		config:         config,
		env:            env,
		db:             db,
		exportImportDB: exportImportDB,
		publishDB:      publishDB,
	}, nil
}

// Routes defines and returns the routes for this server.
func (s *Server) Routes(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/schedule", s.handleSchedule(ctx))
	mux.HandleFunc("/import", s.handleImport(ctx))
	mux.Handle("/health", server.HandleHealthz(ctx))

	return mux
}
