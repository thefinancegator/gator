package handlers

import (
	"fmt"

	"github.com/thefinancegator/gator/internal/config"
	"github.com/thefinancegator/gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Handler func(s *State, cmd Command) error

type Commands struct {
	registry map[string]Handler
}

func NewCommands() *Commands {
	return &Commands{registry: map[string]Handler{}}
}

func (c *Commands) Register(name string, f Handler) {
	c.registry[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.registry[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(s, cmd)
}
