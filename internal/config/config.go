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

// Package config defines how the client is configured.
package config

import (
	"errors"
	"fmt"
	"os"

	gap "github.com/muesli/go-app-paths"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultFile = "go-rest.yaml"
	DefaultPath = "$XDG_CONFIG_HOME/go-rest/" + defaultFile
)

var (
	AppScope = gap.NewScope(gap.User, "go-rest")

	config = new(Config)
)

// Config is the configuration for the client.
type Config struct {
	// DataDir is the directory where request and environment data is stored.
	DataDir string `json:"data_dir,omitempty" mapstructure:"data_dir"`
	// DefaultEnv is the name of the environment to use by default whenever a new session is started.
	DefaultEnv string `json:"default_env,omitempty" mapstructure:"default_env"`
	// ColorScheme sets which color scheme will be used.
	ColorScheme string `json:"color_scheme,omitempty" mapstructure:"color_scheme"`
	// Log is the configuration for logging.
	Log Log `json:"log,omitempty" mapstructure:"log"`
}

// Log contains all configuration options for logging.
type Log struct {
	// Level is the level of logging to use.
	Level string `json:"level,omitempty" mapstructure:"level"`
	// Path is the full filepath of the file for logs to be written to. By default, this is a temporary file in $HOME/.config/go-rest/log.
	Path string `json:"path,omitempty" mapstructure:"path"`
	// Format determines which format logs are written in. Supported values are: "json", "text".
	Format string `json:"format,omitempty" mapstructure:"format"`
}

// Load reads the file at configFile and parses it.
func Load() error {
	err := viper.BindPFlag("config", flag.Lookup("config"))
	if err != nil {
		return fmt.Errorf("binding flag: %w", err)
	}

	if err = setDefaults(); err != nil {
		return fmt.Errorf("setting defaults: %w", err)
	}

	configFile := os.ExpandEnv(viper.GetString("config"))
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	config.DataDir = os.ExpandEnv(config.DataDir)
	config.Log.Path = os.ExpandEnv(config.Log.Path)

	return nil
}

func setDefaults() error {
	// color scheme
	viper.SetDefault("color_scheme", "default")

	// config file
	filepath, err := AppScope.ConfigPath(defaultFile)
	if err != nil {
		return fmt.Errorf("default file: %w", err)
	}
	viper.SetDefault("config", filepath)

	// data directory
	dirs, err := AppScope.DataDirs()
	if err != nil {
		return fmt.Errorf("finding app data directories: %w", err)
	}

	if len(dirs) == 0 {
		return errors.New("unable to determine default app data directory")
	}
	viper.SetDefault("data_dir", dirs[0])

	// log path
	logPath, err := AppScope.LogPath("go-rest.log")
	if err != nil {
		return fmt.Errorf("log path: %w", err)
	}
	viper.SetDefault("log.path", logPath)
	// log level
	viper.SetDefault("log.level", "info")
	// log format
	viper.SetDefault("log.format", "json")

	return nil
}

func Logging() Log {
	return config.Log
}

func DataDir() string {
	return config.DataDir
}

func DefaultEnv() string {
	return config.DefaultEnv
}

func ColorScheme() string {
	return config.ColorScheme
}
