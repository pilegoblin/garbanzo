-- name: GetPostsFromChannel :many
SELECT
    *
FROM
    posts
WHERE
    channel_id = $1
ORDER BY
    created_at;

-- name: GetUserWithEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: GetChannel :one
SELECT
    *
FROM
    channels
WHERE
    id = $1;

-- name: GetChannelsInPod :many
SELECT
    *
FROM
    channels
WHERE
    pod_id = $1;

-- name: GetPods :many
SELECT
    pod_id
FROM
    pod_users
WHERE
    user_id = $1;

-- INSERTS:
-- name: CreateChannel :one
INSERT INTO channels(channel_name, pod_id)
    VALUES ($1, $2)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users(email, picture, username)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: CreatePost :one
INSERT INTO posts(content, user_id, channel_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: JoinPod :one
INSERT INTO pod_users(user_id, pod_id)
    VALUES ($1, $2)
RETURNING
    *;

-- name: CreatePod :one
INSERT INTO pods(pod_name, pod_owner)
    VALUES ($1, $2)
RETURNING
    *;

