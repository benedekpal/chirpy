-- name: SetChirpyRedStatus :one
UPDATE users
SET updated_at = NOW(), is_chirpy_red = $1
WHERE id = $2
RETURNING id, created_at, updated_at, email, is_chirpy_red;