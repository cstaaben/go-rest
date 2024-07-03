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

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"
	"log/slog"

	"github.com/cstaaben/go-rest/internal/config"
	"github.com/cstaaben/go-rest/internal/model"
)

func init() {
	flag.StringP("config", "c", config.DefaultPath, "Path to the configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := config.Load()
	if err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}

	closeFn, err := setupLogger(config.Logging())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closeFn()

	slog.Debug("Starting client", slog.String("colorscheme", config.ColorScheme()))

	p := tea.NewProgram(model.New(), tea.WithAltScreen(), tea.WithContext(ctx))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupLogger(cfg config.Log) (func(), error) {
	logFile, err := openLogFile(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("log file: %w", err)
	}

	closeFn := func() {
		logFile.Close() // nolint:errcheck

		info, err := logFile.Stat()
		// only remove the log file if it is empty and the log level is not debug
		if err == nil && info.Size() == 0 && !strings.EqualFold(cfg.Level, "debug") {
			os.Remove(logFile.Name()) // nolint:errcheck
		}
	}

	lvl := new(slog.LevelVar)
	err = lvl.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return closeFn, fmt.Errorf("log level: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level:     lvl.Level(),
		AddSource: true,
	}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(logFile, opts)
	case "text":
		handler = slog.NewTextHandler(logFile, opts)
	}

	slog.SetDefault(slog.New(handler))

	return closeFn, nil
}

func openLogFile(filepath string) (*os.File, error) {
	var (
		file *os.File
		err  error
	)

	// use the default log directory
	if filepath == "" {
		// check if the log directory exists
		_, err = os.Stat(path.Join(config.DefaultPath, "log"))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("log directory: %w", err)
		} else if err != nil && errors.Is(err, os.ErrNotExist) {
			// create the log directory
			err = os.Mkdir(path.Join(config.DefaultPath, "log"), 0755)
			if err != nil {
				return nil, fmt.Errorf("log directory: %w", err)
			}
		}

		// create the log file
		file, err = os.CreateTemp(
			os.ExpandEnv(config.Logging().Path),
			fmt.Sprintf("go-rest.%s.*.log", time.Now().Format("20060102")),
		)
	} else {
		file, err = os.OpenFile(os.ExpandEnv(filepath), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		return nil, fmt.Errorf("log file: %w", err)
	}

	return file, nil
}
