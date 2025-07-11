--
-- User Queries
--

-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  auth_id,
  avatar_url,
  user_color
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByAuthID :one
SELECT * FROM users
WHERE auth_id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  username = sqlc.arg(username),
  avatar_url = sqlc.arg(avatar_url)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

--
-- Pod Queries
--

-- name: CreatePod :one
INSERT INTO pods (
  owner_id,
  name,
  invite_code
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPod :one
SELECT * FROM pods
WHERE id = $1 LIMIT 1;

-- name: GetPodByInviteCode :one
SELECT * FROM pods
WHERE invite_code = $1 LIMIT 1;

-- name: ListPodsForUser :many
SELECT p.* FROM pods p
JOIN pod_members pm ON p.id = pm.pod_id
WHERE pm.user_id = $1;

-- name: UpdatePod :one
UPDATE pods
SET
  name = sqlc.arg(name),
  invite_code = sqlc.arg(invite_code)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeletePod :exec
DELETE FROM pods
WHERE id = $1;

--
-- Pod Member Queries
--

-- name: AddPodMember :one
INSERT INTO pod_members (
  user_id,
  pod_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetPodMember :one
SELECT * FROM pod_members
WHERE user_id = $1 AND pod_id = $2;

-- name: ListPodMembers :many
SELECT * FROM pod_members
WHERE pod_id = $1;

-- name: RemovePodMember :exec
DELETE FROM pod_members
WHERE user_id = $1 AND pod_id = $2;

-- name: CheckUserInPod :one
SELECT EXISTS (
  SELECT 1 FROM pod_members
  WHERE user_id = $1 AND pod_id = $2
);

--
-- Bean Queries
--

-- name: CreateBean :one
INSERT INTO beans (
  pod_id,
  name
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetBean :one
SELECT * FROM beans
WHERE id = $1 LIMIT 1;

-- name: ListBeansForPod :many
SELECT * FROM beans
WHERE pod_id = $1
ORDER BY name;

-- name: ListBeansForPodFull :many
SELECT
    b.id,
    b.name,
    p.id as pod_id,
    p.name as pod_name,
    COALESCE(jsonb_agg(
        json_build_object(
            'id', m.id,
            'content', m.content,
            'created_at', m.created_at,
            'author_id', u.id,
            'author_username', u.username,
            'author_avatar_url', u.avatar_url,
            'author_user_color', u.user_color
        ) ORDER BY m.created_at ASC
    ) FILTER (WHERE m.id IS NOT NULL), '[]'::jsonb) as messages
FROM beans b
JOIN pods p ON b.pod_id = p.id
LEFT JOIN messages m ON m.bean_id = b.id
LEFT JOIN users u ON m.author_id = u.id
WHERE b.pod_id = $1
GROUP BY b.id, p.id, p.name
ORDER BY b.name;

-- name: UpdateBean :one
UPDATE beans
SET
  name = sqlc.arg(name)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteBean :exec
DELETE FROM beans
WHERE id = $1;

--
-- Message Queries
--

-- name: CreateMessage :one
WITH new_message AS (
  INSERT INTO messages (
    id,
    bean_id,
    author_id,
    content
  ) VALUES (
    $1, $2, $3, $4
  ) RETURNING *
)
SELECT
  m.*,
  u.username as author_username,
  u.avatar_url as author_avatar_url,
  u.user_color as author_user_color,
  u.id as author_id,
  b.pod_id as pod_id,
  b.name as bean_name
FROM new_message m
JOIN users u ON m.author_id = u.id
JOIN beans b ON m.bean_id = b.id;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: ListMessagesInBean :many
SELECT m.*, u.username as author_username, u.id as author_id FROM messages m
JOIN users u ON m.author_id = u.id
WHERE bean_id = $1
ORDER BY m.created_at
LIMIT 50;

-- name: UpdateMessage :one
WITH updated_message AS (
  UPDATE messages
  SET
    content = $2,
    updated_at = now()
  WHERE messages.id = $1
  RETURNING *
)
SELECT
  um.id,
  um.bean_id,
  um.author_id,
  um.content,
  um.created_at,
  um.updated_at,
  b.pod_id,
  b.name as bean_name,
  u.username as author_username,
  u.avatar_url as author_avatar_url,
  u.user_color as author_user_color,
  u.id as author_id
FROM updated_message um
JOIN beans b on um.bean_id = b.id
JOIN users u ON um.author_id = u.id;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;
