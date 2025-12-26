-- name: CreateUserProfilePic :one
INSERT INTO user_images (file_name, file_size, file_type, user_id)
VALUES ($1, $2, $3, $4)
RETURNING file_id, file_name, file_size, file_type, file_uploaded_at, user_id;