-- name: CreateUser :one
INSERT INTO users (user_full_name, user_email, user_google_id, user_profile_pic_id, user_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, user_full_name, user_email, user_google_id, user_profile_pic_id, user_created_at, user_updated_at;

-- name: GetUserByID :one
SELECT 
    u.user_id, 
    u.user_full_name, 
    u.user_email, 
    u.user_profile_pic_id, 
    i.file_name AS user_profile_pic_name,
    u.user_created_at, 
    u.user_updated_at
FROM users u
LEFT JOIN user_images i
    ON u.user_profile_pic_id = i.file_id
WHERE u.user_id = $1;

-- name: GetUserByGoogleId :one
SELECT 
    u.user_id, 
    u.user_full_name, 
    u.user_email, 
    u.user_profile_pic_id, 
    i.file_name AS user_profile_pic_name,
    u.user_created_at, 
    u.user_updated_at
FROM users u
LEFT JOIN user_images i
    ON u.user_profile_pic_id = i.file_id
WHERE u.user_google_id = $1;

-- name: GetUserByEmail :one
SELECT 
    u.user_id, 
    u.user_full_name, 
    u.user_email, 
    u.user_profile_pic_id, 
    i.file_name AS user_profile_pic_name,
    u.user_created_at, 
    u.user_updated_at,
    u.user_password
FROM users u
LEFT JOIN user_images i
    ON u.user_profile_pic_id = i.file_id
WHERE u.user_email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  user_full_name = COALESCE(sqlc.narg('user_full_name'), user_full_name),
  user_email = COALESCE(sqlc.narg('user_email'), user_email),
  user_profile_pic_id = COALESCE(sqlc.narg('user_profile_pic_id'), user_profile_pic_id),
  user_password = COALESCE(sqlc.narg('user_password'), user_password),  
  user_google_id = COALESCE(sqlc.narg('user_google_id'), user_google_id)
WHERE user_id = sqlc.arg('user_id')
RETURNING user_id, user_full_name, user_email, user_profile_pic_id, user_created_at, user_updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1
RETURNING user_id, user_full_name, user_email, user_profile_pic_id, user_created_at, user_updated_at;
