-- name: CreateFeed :one
INSERT INTO feeds (
    created_at,
    updated_at,
    name, 
    url,
    user_id,
    id
) VALUES (
    $1, 
    $2,
    $3,
    $4, 
    $5, 
    $6
)
RETURNING *;

-- name: GetFeeds :many 
SELECT 
    feeds.name,
    feeds.url,
    user_id,
    users.name as user_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;


-- name: GetFeedByUrl :one
SELECT
    *
FROM
    feeds
WHERE
    url = $1;


-- name: MarkFeedFetched :one
UPDATE feeds
SET 
    updated_at = NOW(),
    last_fetched_At = NOW() 
WHERE
    url = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT 
    *
FROM 
    feeds
ORDER BY last_fetched_At ASC NULLS FIRST
LIMIT 1;