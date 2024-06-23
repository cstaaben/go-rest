package environments

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cstaaben/go-rest/internal/environment"
)

type Option func(*Model)

func WithDefaultEnvironment(env string) Option {
	return func(m *Model) {
		m.defaultEnv = env
	}
}

type Model struct {
	defaultEnv string
	dataDir    string

	Environments []*environment.Environment
	Selected     *environment.Environment
}

func New(dataDir string, opts ...Option) *Model {
	m := &Model{
		dataDir:      dataDir,
		Environments: make([]*environment.Environment, 0),
	}

	for _, optFunc := range opts {
		optFunc(m)
	}

	return m
}

func loadEnvironments(dataDir string) ([]*environment.Environment, error) {
	panic("not implemented") // TODO: Implement
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (model *Model) Init() tea.Cmd {
	return func() tea.Msg {
		envs, err := loadEnvironments(model.dataDir)
		if err != nil {
			return fmt.Errorf("loading environments: %w", err)
		}

		model.Environments = envs

		if model.defaultEnv != "" {
			for _, env := range model.Environments {
				if strings.EqualFold(env.Name, model.defaultEnv) {
					model.Selected = env
					break
				}
			}
		}

		return nil
	}
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (model *Model) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	panic("not implemented") // TODO: Implement
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (model *Model) View() string {
	panic("not implemented") // TODO: Implement
}
