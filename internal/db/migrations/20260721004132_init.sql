--
-- go-rest - A TUI for a REST client
-- Copyright (C) 2024  Corbin Staaben
--
-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.
--
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.
--
-- You should have received a copy of the GNU General Public License
-- along with this program.  If not, see <https://www.gnu.org/licenses/>.
--
-- +goose Up
CREATE TABLE IF NOT EXISTS `groups`(
  `uuid` TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS requests(
  `uuid` TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  `method` TEXT NOT NULL,
  group_uuid TEXT,
  url TEXT NOT NULL,
  headers TEXT,
  proto TEXT,
  body TEXT,
  FOREIGN KEY(group_uuid) REFERENCES `groups`(`uuid`)
);
CREATE TABLE IF NOT EXISTS request_history(
  id INTEGER PRIMARY KEY,
  request_uuid TEXT NOT NULL,
  url TEXT NOT NULL,
  headers TEXT,
  `method` TEXT NOT NULL,
  proto TEXT,
  body TEXT,
  received_timestamp TEXT DEFAULT CURRENT_TIMESTAMP,
  duration_ms INTEGER NOT NULL,
  FOREIGN KEY(request_uuid) REFERENCES requests(`uuid`)
);
CREATE TABLE IF NOT EXISTS environments(
  `uuid` TEXT PRIMARY KEY,
  name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS environment_variables(
  environment_uuid TEXT,
  name TEXT NOT NULL,
  `value` BLOB,
  PRIMARY KEY(environment_uuid, name),
  FOREIGN KEY(environment_uuid) REFERENCES environments(`uuid`)
);
-- +goose Down
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS request_history;
DROP TABLE IF EXISTS `groups`;
DROP TABLE IF EXISTS environments;
DROP TABLE IF EXISTS environment_variables;
