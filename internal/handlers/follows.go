package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/thefinancegator/gator/internal/database"
)

func Follow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	follow, err := createFeedFollow(s, user.ID, feed.ID)
	if err != nil {
		return err
	}

	fmt.Println("User: " + follow.UserName)
	fmt.Println("Followed feed: " + follow.FeedName)
	return nil
}

func Following(s *State, cmd Command, user database.User) error {
	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get followed feeds: %w", err)
	}

	if len(follows) == 0 {
		fmt.Printf("User %s doesn't follow any feeds\n", user.Name)
		return nil
	}

	fmt.Printf("User %s follows:\n", user.Name)
	for _, f := range follows {
		fmt.Println(" - " + f.FeedName)
	}

	return nil
}

func Unfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	err = s.Db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("%s unfollowed %s\n", user.Name, feed.Name)
	return nil
}

func createFeedFollow(s *State, userID, feedID uuid.UUID) (database.CreateFeedFollowRow, error) {
	now := time.Now().UTC()
	follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feedID,
	})
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("couldn't create feed follow: %w", err)
	}

	return follow, nil
}
