-- --- USER ---
-- name: CreateUser :one
INSERT INTO users (username, password_hash, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- --- FOLLOW ---
-- name: CreateFollow :one
INSERT INTO follows (follower_id, following_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteFollow :exec
DELETE FROM follows 
WHERE follower_id = $1 AND following_id = $2;

-- --- POST & FEED ---
-- name: CreatePost :one
INSERT INTO posts (user_id, content)
VALUES ($1, $2)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1 LIMIT 1;

-- name: ListPostsByUser :many
SELECT * FROM posts 
WHERE user_id = $1 
ORDER BY created_at DESC;

-- name: GetNewsFeed :many
SELECT 
    posts.id, 
    posts.content, 
    posts.created_at, 
    users.username AS author_name -- JOIN để lấy tên tác giả
FROM posts
JOIN users ON posts.user_id = users.id
WHERE posts.user_id IN (
    SELECT following_id 
    FROM follows 
    WHERE follower_id = $1 -- $1 là ID của chính mình
)
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;