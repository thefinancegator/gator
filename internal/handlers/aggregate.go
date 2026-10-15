package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"

	"github.com/thefinancegator/gator/internal/database"
	"github.com/thefinancegator/gator/internal/requests"
)

var pubDateLayouts = []string{time.RFC1123Z, time.RFC1123, time.RFC3339}

func Aggregate(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse time between requests: %w", err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("time between requests must be positive, got %s", timeBetweenRequests)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			log.Printf("couldn't scrape feeds: %v", err)
		}
	}
}

func scrapeFeeds(s *State) error {
	ctx := context.Background()

	next, err := s.Db.GetNextFeedToFetch(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("no feeds to fetch, add one with addfeed")
	}
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	feed, err := s.Db.MarkFeedFetched(ctx, next.Url)
	if err != nil {
		return fmt.Errorf("couldn't mark feed %s as fetched: %w", next.Name, err)
	}

	rss, err := requests.FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed %s: %w", feed.Name, err)
	}

	for _, item := range rss.Channel.Item {
		now := time.Now().UTC()
		_, err := s.Db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			Url:         item.Link,
			PublishedAt: parsePubDate(item.PubDate),
		})
		if pq.As(err, pqerror.UniqueViolation) != nil {
			continue
		}
		if err != nil {
			log.Printf("couldn't create post %q: %v", item.Title, err)
		}
	}

	log.Printf("Feed %s collected, %d posts found", feed.Name, len(rss.Channel.Item))
	return nil
}

func parsePubDate(value string) sql.NullTime {
	for _, layout := range pubDateLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}
	return sql.NullTime{}
}
