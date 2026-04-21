-- name: CreateUser :one
INSERT INTO users (user_full_name, user_email, user_google_id, user_password)
VALUES ($1, $2, $3, $4)
RETURNING user_id, user_full_name, user_email, user_google_id, user_created_at, user_updated_at;

-- name: GetUserByID :one
SELECT
    u.user_id,
    u.user_full_name,
    u.user_email,
    u.user_created_at,
    u.user_updated_at,
    i.file_name,
    u.is_temp,
    u.is_deleted
FROM users u 
LEFT JOIN user_images i
    ON u.user_id = i.user_id
WHERE u.user_id = $1 AND u.is_deleted = FALSE;

-- name: GetUserPassword :one
SELECT
    u.user_password
FROM users u
WHERE u.user_id = $1 AND u.is_deleted = FALSE AND u.is_temp = FALSE;

-- name: GetUserByGoogleId :one
SELECT
    u.user_id,
    u.user_full_name,
    u.user_email,
    u.user_created_at,
    u.user_updated_at,
    i.file_name
FROM users u
LEFT JOIN user_images i
    ON u.user_id = i.user_id
WHERE u.user_google_id = $1 AND u.is_deleted = FALSE AND u.is_temp = FALSE;

-- name: GetUserByEmail :one
SELECT
    u.user_id,
    u.user_full_name,
    u.user_email,
    u.user_created_at,
    u.user_updated_at,
    u.user_password,
    i.file_name
FROM users u
LEFT JOIN user_images i
    ON u.user_id = i.user_id
WHERE u.user_email = $1 AND u.is_deleted = FALSE AND u.is_temp = FALSE;

-- name: UpdateUser :one
UPDATE users
SET
  user_full_name = COALESCE(sqlc.narg('user_full_name'), user_full_name),
  user_email = COALESCE(sqlc.narg('user_email'), user_email),
  user_password = COALESCE(sqlc.narg('user_password'), user_password),
  user_google_id = COALESCE(sqlc.narg('user_google_id'), user_google_id)
WHERE user_id = sqlc.arg('user_id') AND is_deleted = FALSE AND is_temp = FALSE
RETURNING user_id, user_full_name, user_email, user_created_at, user_updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1
RETURNING user_id, user_full_name, user_email, user_created_at, user_updated_at;

-- name: CreateUserProfilePic :one
INSERT INTO user_images (file_name, file_size, file_type, user_id)
VALUES ($1, $2, $3, $4)
RETURNING file_id, file_name, file_size, file_type, user_id;

-- name: UpdateUserProfilePic :one
UPDATE user_images
SET
  file_name = COALESCE(sqlc.narg(file_name), file_name),
  file_size = COALESCE(sqlc.narg(file_size), file_size),
  file_type = COALESCE(sqlc.narg(file_type), file_type),
  file_uploaded_at = COALESCE(sqlc.narg(file_uploaded_at), file_uploaded_at)
WHERE user_id = sqlc.arg('user_id')
RETURNING file_id, file_name, file_size, file_type, user_id;

-- name: DeleteUserProfilePic :exec
DELETE FROM user_images
WHERE user_id = $1
RETURNING file_id, file_name, file_size, file_type, user_id;

-- name: GetUserProfilePic :one
SELECT *
FROM user_images
WHERE user_id = $1 AND is_deleted = FALSE AND is_temp = FALSE;

-- name: CreateTempUser :one
INSERT INTO users (is_temp, user_full_name)
VALUES (TRUE, sqlc.arg('user_full_name'))
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET is_deleted = TRUE
WHERE user_id = sqlc.arg('user_id');