-- name: CreateForm :one
INSERT INTO forms (form_title, form_description, created_by)
VALUES ($1, $2, $3)
RETURNING form_id, form_title, form_description, created_by, form_created_at, form_updated_at, form_status, form_access;

-- name: UpdateForm :one
UPDATE forms
SET
  form_title = COALESCE(sqlc.narg('form_title'), form_title),
  form_description = COALESCE(sqlc.narg('form_description'), form_description),
  created_by = COALESCE(sqlc.narg('created_by'), created_by),
  form_status = COALESCE(sqlc.narg('form_status'), form_status),
  form_access = COALESCE(sqlc.narg('form_access'), form_access)
WHERE form_id = sqlc.arg('form_id')
RETURNING form_id, form_title, form_description, created_by, form_created_at, form_updated_at, form_status, form_access;

-- name: GetFormByID :one
SELECT *
FROM forms
WHERE form_id = $1;

-- name: ListForms :many
SELECT *
FROM forms
WHERE
    created_by = sqlc.arg('created_by')
    AND (sqlc.narg('search') IS NULL
         OR form_title ILIKE '%' || sqlc.narg('search') || '%'
         OR form_description ILIKE '%' || sqlc.narg('search') || '%')
    AND (sqlc.narg('form_status') IS NULL OR form_status = sqlc.narg('form_status'))
    AND (sqlc.narg('form_access') IS NULL OR form_access = sqlc.narg('form_access'))
ORDER BY
    CASE LOWER(COALESCE(sqlc.narg('short_by'), 'created'))
        WHEN 'title' THEN form_title
        WHEN 'status' THEN form_status
        WHEN 'access' THEN form_access
        WHEN 'updated' THEN form_updated_at
        ELSE form_created_at
    END
    /* Direction: 'asc' or 'desc' */
    COLLATE "C" 
    /* Note: SQLC currently does not directly allow dynamic ASC/DESC, handle in app code if needed */
LIMIT COALESCE(sqlc.narg('limit'), 10)
OFFSET COALESCE(sqlc.narg('offset'), 0);


-- name: DeleteForm :one
DELETE FROM forms
WHERE form_id = $1
RETURNING form_id, form_title, form_description, created_by, form_created_at, form_updated_at, form_status, form_access;
