--
-- User Queries
--

-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  auth_id,
  avatar_url
) VALUES (
  $1, $2, $3, $4
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
INSERT INTO messages (
  bean_id,
  author_id,
  content
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: ListMessagesInBean :many
SELECT m.*, u.username as author_username, u.avatar_url as author_avatar_url FROM messages m
JOIN users u ON m.author_id = u.id
WHERE bean_id = $1
ORDER BY m.created_at
LIMIT 50;

-- name: UpdateMessage :one
UPDATE messages
SET
  content = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;
