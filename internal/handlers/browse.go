package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/thefinancegator/gator/internal/database"
)

func Browse(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	}

	limit := int32(2)
	if len(cmd.Args) == 1 {
		n, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil || n < 1 {
			return fmt.Errorf("invalid limit %q: must be a positive integer", cmd.Args[0])
		}
		limit = int32(n)
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		published := "unknown date"
		if post.PublishedAt.Valid {
			published = post.PublishedAt.Time.Format("Mon Jan 2")
		}
		fmt.Printf("%s from %s\n", published, post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
