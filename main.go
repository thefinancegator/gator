package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/pawelwerner/gator/internal/config"
	"github.com/pawelwerner/gator/internal/database"
	"github.com/pawelwerner/gator/internal/handlers"
	"github.com/pawelwerner/gator/internal/middleware"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gator <command> [args...]")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	state := &handlers.State{
		Db:     database.New(db),
		Config: &cfg,
	}

	cmds := handlers.NewCommands()
	cmds.Register("login", handlers.Login)
	cmds.Register("register", handlers.Register)
	cmds.Register("reset", handlers.Reset)
	cmds.Register("users", handlers.ListUsers)
	cmds.Register("agg", handlers.Aggregate)
	cmds.Register("addfeed", middleware.LoggedIn(handlers.AddFeed))
	cmds.Register("feeds", handlers.ListFeeds)
	cmds.Register("follow", middleware.LoggedIn(handlers.Follow))
	cmds.Register("following", middleware.LoggedIn(handlers.Following))
	cmds.Register("unfollow", middleware.LoggedIn(handlers.Unfollow))
	cmds.Register("browse", middleware.LoggedIn(handlers.Browse))

	cmd := handlers.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cmds.Run(state, cmd); err != nil {
		log.Fatal(err)
	}
}
