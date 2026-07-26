/*
 * go-rest - A TUI for a REST client
 * Copyright (C) 2024  Corbin Staaben
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package database contains database logic for retrieving and updating data.
package database

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/cstaaben/go-rest/internal/config"
	"github.com/cstaaben/go-rest/internal/database/dto"
	"github.com/cstaaben/go-rest/internal/slogutils"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var ErrNotFound = errors.New("not found")

type DAO struct {
	db *sqlx.DB

	requestListDataStmt      *sqlx.Stmt
	requestStmt              *sqlx.Stmt
	groupRequestsStmt        *sqlx.Stmt
	groupsStmt               *sqlx.Stmt
	environmentsStmt         *sqlx.Stmt
	environmentStmt          *sqlx.Stmt
	environmentVariablesStmt *sqlx.Stmt
	environmentVariableStmt  *sqlx.Stmt
}

// New initializes and returns a new DAO instance, connecting to SQLite and running migrations.
func New(ctx context.Context) (*DAO, error) {
	db := &DAO{}

	var (
		dbName string
		err    error

		dir        = config.DataDir()
		devEnabled = viper.GetBool("dev")
	)
	if devEnabled {
		dbName = "development"
	} else {
		dbName = "data"
	}

	db.db, err = sqlx.ConnectContext(ctx, "sqlite", filepath.Join(dir, fmt.Sprintf("%s.db", dbName)))
	if err != nil {
		slog.Error("Failed to connect to database", slogutils.Err(err))
		return nil, err
	}

	if _, err = db.db.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
		slog.Error("Failed to enable foreign keys", slogutils.Err(err))
		return nil, err
	}

	if err = db.runMigrations(ctx); err != nil {
		slog.Error("Failed to setup database", slogutils.Err(err))
		return nil, err
	}

	if err := db.prepareQueries(ctx); err != nil {
		slog.Error("Failed to prepare database", slogutils.Err(err))
		return nil, err
	}

	return db, nil
}

// TODO: check version before running migrations?
func (dao *DAO) runMigrations(ctx context.Context) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("Failed to set migration dialect", slogutils.Err(err))
		return err
	}

	if err := goose.UpContext(ctx, dao.db.DB, "migrations"); err != nil {
		slog.Error("Failed to migrate database", slogutils.Err(err))
		return err
	}

	return nil
}

func (dao *DAO) prepareQueries(ctx context.Context) error {
	var err error

	requestListDataQuery := `SELECT
    r.uuid AS request_id,
    r.name AS request_name,
    r.description AS request_desc,
    r.method AS request_method,
    g.uuid AS group_id,
    g.name AS group_name,
    g.description AS group_desc
FROM requests r
JOIN groups g ON r.group_uuid = g.uuid
WHERE g.uuid = 'unsorted'

UNION ALL

SELECT
    '' AS request_id,
    '' AS request_name,
    NULL AS request_desc,
    '' AS request_method,
    g.uuid AS group_id,
    g.name AS group_name,
    g.description AS group_desc
FROM groups g
WHERE g.uuid != 'unsorted'
ORDER BY group_name ASC, request_name ASC
LIMIT ? OFFSET ?;`
	dao.requestListDataStmt, err = dao.db.PreparexContext(ctx, requestListDataQuery)
	if err != nil {
		return fmt.Errorf("request list data: %w", err)
	}

	requestQuery := `SELECT uuid, name, description AS desc, method, group_uuid, url, headers, proto, body FROM requests WHERE uuid = ?;`
	dao.requestStmt, err = dao.db.PreparexContext(ctx, requestQuery)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	groupRequestsQuery := `SELECT uuid, name, description AS desc, method, group_uuid, url, headers, proto, body FROM requests WHERE group_uuid = ?;`
	dao.groupRequestsStmt, err = dao.db.PreparexContext(ctx, groupRequestsQuery)
	if err != nil {
		return fmt.Errorf("group requests: %w", err)
	}

	groupsQuery := `SELECT uuid, name, description AS desc FROM groups ORDER BY uuid ASC;`
	dao.groupsStmt, err = dao.db.PreparexContext(ctx, groupsQuery)
	if err != nil {
		return fmt.Errorf("groups: %w", err)
	}

	environmentsQuery := `SELECT uuid, name FROM environments ORDER BY name ASC, uuid ASC;`
	dao.environmentsStmt, err = dao.db.PreparexContext(ctx, environmentsQuery)
	if err != nil {
		return fmt.Errorf("environments: %w", err)
	}

	environmentQuery := `SELECT uuid, name FROM environments WHERE uuid = ?;`
	dao.environmentStmt, err = dao.db.PreparexContext(ctx, environmentQuery)
	if err != nil {
		return fmt.Errorf("environment: %w", err)
	}

	environmentVariablesQuery := `SELECT environment_uuid, name, value FROM environment_variables WHERE environment_uuid = ? ORDER BY name ASC;`
	dao.environmentVariablesStmt, err = dao.db.PreparexContext(ctx, environmentVariablesQuery)
	if err != nil {
		return fmt.Errorf("environment variables: %w", err)
	}

	environmentVariableQuery := `SELECT environment_uuid, name, value FROM environment_variables WHERE environment_uuid = ? AND name = ?;`
	dao.environmentVariableStmt, err = dao.db.PreparexContext(ctx, environmentVariableQuery)
	if err != nil {
		return fmt.Errorf("environment variable: %w", err)
	}

	return nil
}

// Close closes the underlying database connection.
func (dao *DAO) Close() error {
	return dao.db.Close()
}

// RequestListData retrieves a paginated slice of requests and associated group metadata.
// Parameters limit and offset specify the maximum number of items and the query offset.
func (dao *DAO) RequestListData(ctx context.Context, limit, offset int) ([]*dto.RequestListData, error) {
	dest := make([]*dto.RequestListData, 0)
	err := dao.requestListDataStmt.SelectContext(ctx, &dest, limit, offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get request list data", slogutils.Err(err))
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("No requests found")
		return dest, nil
	}

	return dest, nil
}

type requestRow struct {
	UUID      string `db:"uuid"`
	Name      string `db:"name"`
	Desc      string `db:"desc"`
	Method    string `db:"method"`
	GroupUUID string `db:"group_uuid"`
	URL       string `db:"url"`
	Headers   string `db:"headers"`
	Proto     string `db:"proto"`
	Body      string `db:"body"`
}

func (row *requestRow) toDTO() *dto.Request {
	var headers map[string][]string
	if row.Headers != "" {
		_ = json.Unmarshal([]byte(row.Headers), &headers)
	}
	return &dto.Request{
		UUID:      row.UUID,
		Name:      row.Name,
		Desc:      row.Desc,
		GroupUUID: row.GroupUUID,
		Data: &dto.Data{
			URL:     row.URL,
			Headers: headers,
			Method:  row.Method,
			Proto:   row.Proto,
			Body:    row.Body,
		},
	}
}

// Request retrieves a request by its UUID.
// Returns ErrNotFound if no request with the given ID exists.
func (dao *DAO) Request(ctx context.Context, requestID string) (*dto.Request, error) {
	var row requestRow
	err := dao.requestStmt.GetContext(ctx, &row, requestID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get request data", slogutils.Err(err), slog.String("request_id", requestID))
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("No request found for provided ID", slog.String("request_id", requestID))
		return nil, ErrNotFound
	}

	return row.toDTO(), nil
}

// GroupRequests retrieves all requests belonging to a specific group UUID.
// Returns ErrNotFound if the group contains no requests or does not exist.
func (dao *DAO) GroupRequests(ctx context.Context, groupID string) ([]*dto.Request, error) {
	var rows []requestRow
	err := dao.groupRequestsStmt.SelectContext(ctx, &rows, groupID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get group requests", slogutils.Err(err), slog.String("group_id", groupID))
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) || len(rows) == 0 {
		slog.Debug("No requests found for provided group", slog.String("group_id", groupID))
		return nil, ErrNotFound
	}

	res := make([]*dto.Request, len(rows))
	for i, r := range rows {
		res[i] = r.toDTO()
	}
	return res, nil
}

// Groups retrieves all request groups sorted by UUID in ascending order.
func (dao *DAO) Groups(ctx context.Context) ([]*dto.Group, error) {
	groups := make([]*dto.Group, 0)
	err := dao.groupsStmt.SelectContext(ctx, &groups)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get request groups", slogutils.Err(err))
		return nil, err
	}

	return groups, nil
}

// UpdateRequest updates an existing request's attributes in the database.
func (dao *DAO) UpdateRequest(ctx context.Context, req *dto.Request) error {
	query := `UPDATE requests 
SET
	name = ?,
	description = ?,
	method = ?,
	group_uuid = ?,
	url = ?,
	headers = ?,
	proto = ?,
	body = ?
WHERE uuid = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to start transaction to update request", slogutils.Err(err), slog.String("request_name", req.Name), slog.String("request_id", req.UUID))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	headerBytes, err := json.Marshal(req.Data.Headers)
	if err != nil {
		slog.Error("Failed to encode request headers", slogutils.Err(err), slog.Any("headers", req.Data.Headers))
		return err
	}

	_, err = tx.ExecContext(ctx, query, req.Name, req.Desc, req.Data.Method, req.GroupUUID, req.Data.URL, string(headerBytes), req.Data.Proto, req.Data.Body, req.UUID)
	if err != nil {
		slog.Error("Failed to update request", slogutils.Err(err), slog.String("request_id", req.UUID), slog.String("request_name", req.Name))
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("Failed to commit request update", slogutils.Err(err), slog.String("request_id", req.UUID), slog.String("request_name", req.Name))
		return err
	}

	return nil
}

// UpdateGroup updates an existing request group's name and description in the database.
func (dao *DAO) UpdateGroup(ctx context.Context, group *dto.Group) error {
	logger := slog.With(slog.String("group_id", group.UUID))
	query := `UPDATE groups
SET
	name = ?,
	description = ?
WHERE uuid = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for group update", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	_, err = tx.ExecContext(ctx, query, group.Name, group.Desc, group.UUID)
	if err != nil {
		logger.Error("Failed to update request group", slogutils.Err(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("Failed to commit group update", slogutils.Err(err))
		return err
	}

	return nil
}

// InsertRequest inserts a new request into the database.
func (dao *DAO) InsertRequest(ctx context.Context, request *dto.Request) error {
	query := `INSERT INTO requests (
	uuid,
	name,
	description,
	method,
	group_uuid,
	url,
	headers,
	proto,
	body
) VALUES (
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?
);`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to start transaction to insert new request", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	headerBytes, err := json.Marshal(request.Data.Headers)
	if err != nil {
		slog.Error("Failed to encode request headers", slogutils.Err(err), slog.Any("headers", request.Data.Headers))
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		request.UUID,
		request.Name,
		request.Desc,
		request.Data.Method,
		request.GroupUUID,
		request.Data.URL,
		string(headerBytes),
		request.Data.Proto,
		request.Data.Body,
	)
	if err != nil {
		slog.Error("Failed to insert new request", slogutils.Err(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("Failed to commit new request", slogutils.Err(err))
		return err
	}

	return nil
}

// InsertGroup inserts a new request group into the database.
func (dao *DAO) InsertGroup(ctx context.Context, group *dto.Group) error {
	logger := slog.With(slog.String("group_id", group.UUID), slog.String("group_name", group.Name))
	query := `INSERT INTO groups (
	uuid,
	name,
	description
) VALUES (
	?,
	?,
	?
);`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for group insert", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	_, err = tx.ExecContext(ctx, query, group.UUID, group.Name, group.Desc)
	if err != nil {
		logger.Error("Failed to insert request group", slogutils.Err(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit group insert", slogutils.Err(err))
		return err
	}

	return nil
}

// DeleteRequest deletes a request by its UUID from the database.
// Returns ErrNotFound if no request with the given ID exists.
func (dao *DAO) DeleteRequest(ctx context.Context, requestID string) error {
	logger := slog.With(slog.String("request_id", requestID))
	query := `DELETE FROM requests WHERE uuid = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for request deletion", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	result, err := tx.ExecContext(ctx, query, requestID)
	if err != nil {
		logger.Error("Failed to delete request", slogutils.Err(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		logger.Debug("No request found to delete", slog.String("request_id", requestID))
		return ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit request deletion", slogutils.Err(err))
		return err
	}

	return nil
}

// DeleteGroup deletes a request group by its UUID from the database.
// Deleting the default 'unsorted' group is not permitted and returns a wrapped [errors.ErrUnsupported] error.
// Any requests associated with the deleted group are automatically deleted via database cascade.
// Returns ErrNotFound if no group with the given ID exists.
func (dao *DAO) DeleteGroup(ctx context.Context, groupID string) error {
	if groupID == dto.UnsortedName {
		return fmt.Errorf("%w: cannot delete unsorted group", errors.ErrUnsupported)
	}

	logger := slog.With(slog.String("group_id", groupID))

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for group deletion", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	deleteQuery := `DELETE FROM groups WHERE uuid = ?;`
	result, err := tx.ExecContext(ctx, deleteQuery, groupID)
	if err != nil {
		logger.Error("Failed to delete request group", slogutils.Err(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		logger.Debug("No group found to delete", slog.String("group_id", groupID))
		return ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit group deletion", slogutils.Err(err))
		return err
	}

	return nil
}

// InsertEnvironment inserts a new environment and any contained variables into the database.
func (dao *DAO) InsertEnvironment(ctx context.Context, env *dto.Environment) error {
	logger := slog.With(slog.String("environment_id", env.UUID), slog.String("environment_name", env.Name))
	query := `INSERT INTO environments (uuid, name) VALUES (?, ?);`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for environment insert", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	_, err = tx.ExecContext(ctx, query, env.UUID, env.Name)
	if err != nil {
		logger.Error("Failed to insert environment", slogutils.Err(err))
		return err
	}

	varQuery := `INSERT INTO environment_variables (environment_uuid, name, value) VALUES (?, ?, ?);`
	for _, v := range env.Variables {
		_, err = tx.ExecContext(ctx, varQuery, env.UUID, v.Name, v.Value)
		if err != nil {
			logger.Error("Failed to insert environment variable", slogutils.Err(err), slog.String("variable_name", v.Name))
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit environment insert", slogutils.Err(err))
		return err
	}

	return nil
}

// Environment retrieves an environment by its UUID, including all associated variables.
// Returns ErrNotFound if no environment with the given ID exists.
func (dao *DAO) Environment(ctx context.Context, envID string) (*dto.Environment, error) {
	var env dto.Environment
	err := dao.environmentStmt.GetContext(ctx, &env, envID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get environment data", slogutils.Err(err), slog.String("environment_id", envID))
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("No environment found for provided ID", slog.String("environment_id", envID))
		return nil, ErrNotFound
	}

	vars := make([]*dto.EnvironmentVariable, 0)
	err = dao.environmentVariablesStmt.SelectContext(ctx, &vars, envID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get environment variables", slogutils.Err(err), slog.String("environment_id", envID))
		return nil, err
	}
	env.Variables = vars

	return &env, nil
}

// Environments retrieves all environments sorted by name and UUID.
func (dao *DAO) Environments(ctx context.Context) ([]*dto.Environment, error) {
	envs := make([]*dto.Environment, 0)
	err := dao.environmentsStmt.SelectContext(ctx, &envs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get environments", slogutils.Err(err))
		return nil, err
	}

	for _, env := range envs {
		vars := make([]*dto.EnvironmentVariable, 0)
		err = dao.environmentVariablesStmt.SelectContext(ctx, &vars, env.UUID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			slog.Error("Failed to get environment variables", slogutils.Err(err), slog.String("environment_id", env.UUID))
			return nil, err
		}
		env.Variables = vars
	}

	return envs, nil
}

// UpdateEnvironment updates an existing environment's name in the database.
// Returns ErrNotFound if no environment with the given ID exists.
func (dao *DAO) UpdateEnvironment(ctx context.Context, env *dto.Environment) error {
	logger := slog.With(slog.String("environment_id", env.UUID))
	query := `UPDATE environments SET name = ? WHERE uuid = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for environment update", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	result, err := tx.ExecContext(ctx, query, env.Name, env.UUID)
	if err != nil {
		logger.Error("Failed to update environment", slogutils.Err(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		logger.Debug("No environment found to update", slog.String("environment_id", env.UUID))
		return ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit environment update", slogutils.Err(err))
		return err
	}

	return nil
}

// DeleteEnvironment deletes an environment by its UUID from the database.
// Associated environment variables are automatically deleted via database cascade.
// Returns ErrNotFound if no environment with the given ID exists.
func (dao *DAO) DeleteEnvironment(ctx context.Context, envID string) error {
	logger := slog.With(slog.String("environment_id", envID))
	query := `DELETE FROM environments WHERE uuid = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for environment deletion", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	result, err := tx.ExecContext(ctx, query, envID)
	if err != nil {
		logger.Error("Failed to delete environment", slogutils.Err(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		logger.Debug("No environment found to delete", slog.String("environment_id", envID))
		return ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit environment deletion", slogutils.Err(err))
		return err
	}

	return nil
}

// InsertEnvironmentVariable inserts a new environment variable into the database.
func (dao *DAO) InsertEnvironmentVariable(ctx context.Context, envVar *dto.EnvironmentVariable) error {
	logger := slog.With(slog.String("environment_id", envVar.EnvironmentUUID), slog.String("variable_name", envVar.Name))
	query := `INSERT INTO environment_variables (environment_uuid, name, value) VALUES (?, ?, ?);`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for variable insert", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	_, err = tx.ExecContext(ctx, query, envVar.EnvironmentUUID, envVar.Name, envVar.Value)
	if err != nil {
		logger.Error("Failed to insert environment variable", slogutils.Err(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit environment variable insert", slogutils.Err(err))
		return err
	}

	return nil
}

// EnvironmentVariables retrieves all variables for a given environment UUID.
// Returns ErrNotFound if no variables exist for the environment.
func (dao *DAO) EnvironmentVariables(ctx context.Context, envID string) ([]*dto.EnvironmentVariable, error) {
	vars := make([]*dto.EnvironmentVariable, 0)
	err := dao.environmentVariablesStmt.SelectContext(ctx, &vars, envID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get environment variables", slogutils.Err(err), slog.String("environment_id", envID))
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) || len(vars) == 0 {
		slog.Debug("No environment variables found for provided environment ID", slog.String("environment_id", envID))
		return nil, ErrNotFound
	}

	return vars, nil
}

// EnvironmentVariable retrieves a specific environment variable by environment UUID and variable name.
// Returns ErrNotFound if the variable does not exist.
func (dao *DAO) EnvironmentVariable(ctx context.Context, envID, name string) (*dto.EnvironmentVariable, error) {
	var envVar dto.EnvironmentVariable
	err := dao.environmentVariableStmt.GetContext(ctx, &envVar, envID, name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("Failed to get environment variable", slogutils.Err(err), slog.String("environment_id", envID), slog.String("variable_name", name))
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("No environment variable found for provided environment ID and name", slog.String("environment_id", envID), slog.String("variable_name", name))
		return nil, ErrNotFound
	}

	return &envVar, nil
}

// UpdateEnvironmentVariable updates the value of an existing environment variable.
// Returns ErrNotFound if the variable does not exist.
func (dao *DAO) UpdateEnvironmentVariable(ctx context.Context, envVar *dto.EnvironmentVariable) error {
	logger := slog.With(slog.String("environment_id", envVar.EnvironmentUUID), slog.String("variable_name", envVar.Name))
	query := `UPDATE environment_variables SET value = ? WHERE environment_uuid = ? AND name = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for variable update", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	result, err := tx.ExecContext(ctx, query, envVar.Value, envVar.EnvironmentUUID, envVar.Name)
	if err != nil {
		logger.Error("Failed to update environment variable", slogutils.Err(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		logger.Debug("No environment variable found to update", slog.String("environment_id", envVar.EnvironmentUUID), slog.String("variable_name", envVar.Name))
		return ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit environment variable update", slogutils.Err(err))
		return err
	}

	return nil
}

// DeleteEnvironmentVariable deletes a specific environment variable by environment UUID and variable name.
// Returns ErrNotFound if the variable does not exist.
func (dao *DAO) DeleteEnvironmentVariable(ctx context.Context, envID, name string) error {
	logger := slog.With(slog.String("environment_id", envID), slog.String("variable_name", name))
	query := `DELETE FROM environment_variables WHERE environment_uuid = ? AND name = ?;`

	tx, err := dao.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction for variable deletion", slogutils.Err(err))
		return err
	}
	defer tx.Rollback() // nolint:errcheck

	result, err := tx.ExecContext(ctx, query, envID, name)
	if err != nil {
		logger.Error("Failed to delete environment variable", slogutils.Err(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		logger.Debug("No environment variable found to delete", slog.String("environment_id", envID), slog.String("variable_name", name))
		return ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Failed to commit environment variable deletion", slogutils.Err(err))
		return err
	}

	return nil
}
