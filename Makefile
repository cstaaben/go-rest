#
# go-rest - A TUI for a REST client
# Copyright (C) 2024  Corbin Staaben
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#

# ENVIRONMENT is the name of the database environment to use.
ENVIRONMENT ?= production

# DATA_DIR is the path to the directory containing the database files used by go-rest.
DATA_DIR ?= $(HOME)/.local/share/go-rest

# GOOSE_CMD is the standard goose command prefix to use for migrations.
GOOSE_CMD = goose sqlite3 $(DATA_DIR)/$(ENVIRONMENT).db

# Install the binary if no targets are given
all: install

# Build the Go binary
.PHONY: build
build: tidy db-setup
	go build -o bin/go-rest ./cmd/go-rest

# Tidy Go modules
.PHONY: tidy
tidy:
	go mod tidy

# Install the Go binary. Requires $GOBIN to be in $PATH.
.PHONY: install
install:
	go install ./cmd/go-rest

#
# Migrations
#

# Setup the database instances.
.PHONY: db-setup
db-setup:
	$(MAKE) ENVIRONMENT=production migrate-up
	$(MAKE) ENVIRONMENT=development migrate-up

# Apply all migrations to the database.
.PHONY: migrate-up
migrate-up:
	$(GOOSE_CMD) up

# Roll back a single migration version.
.PHONY: migrate-down
migrate-down:
	$(GOOSE_CMD) down

# Create a new SQL migration.
.PHONY: create-sql-migration
create-sql-migration:
	$(GOOSE_CMD) create $(NAME) sql

# Create a new programmatic Go migration.
.PHONY: create-go-migration
create-go-migration:
	$(GOOSE_CMD) create $(NAME)

# Print the current status of all migrations.
.PHONY: migration-status
migration-status:
	$(GOOSE_CMD) status

#
# Convenience targets
#

# Run the binary and dump log file to stdout.
.PHONY: run
run: tidy build
	bin/go-test --config ./example/config.yaml && cat ./logs/test.log

