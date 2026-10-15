package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/thefinancegator/gator/internal/database"
	"github.com/thefinancegator/gator/internal/handlers"
)

func LoggedIn(handler func(s *handlers.State, cmd handlers.Command, user database.User) error) handlers.Handler {
	return func(s *handlers.State, cmd handlers.Command) error {
		if s.Config.CurrentUserName == "" {
			return errors.New("not logged in: run login or register first")
		}

		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't find current user %q: %w", s.Config.CurrentUserName, err)
		}

		return handler(s, cmd, user)
	}
}
