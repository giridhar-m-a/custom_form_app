-- name: CreateForm :one
INSERT INTO forms (
  form_title,
  form_description,
  created_by,
  form_access,
  scheduled_time,
  closing_time,
  is_scheduled,
  invitation_schedule_id,
  invitation_schedule_gap,
  scheduling_id
)
VALUES (
  sqlc.arg('form_title'),
  sqlc.narg('form_description'),
  sqlc.arg('created_by'),
  sqlc.narg('form_access'),
  sqlc.narg('scheduled_time'),
  sqlc.narg('closing_time'),
  sqlc.narg('is_scheduled'),
  sqlc.narg('invitation_schedule_id'),
  sqlc.narg('invitation_schedule_gap'),
  sqlc.narg('scheduling_id')
)
RETURNING *;



-- name: UpdateForm :one
UPDATE forms
SET
  form_title = COALESCE(sqlc.narg('form_title'), form_title),
  form_description = COALESCE(sqlc.narg('form_description'), form_description),
  form_status = COALESCE(sqlc.narg('form_status'), form_status),
  form_access = COALESCE(sqlc.narg('form_access'), form_access),
  scheduling_id = COALESCE(sqlc.narg('scheduling_id'), scheduling_id),
  scheduled_time = COALESCE(sqlc.narg('scheduled_time'), scheduled_time),
  closing_time = COALESCE(sqlc.narg('closing_time'), closing_time),
  is_schedule_completed = COALESCE(sqlc.narg('is_schedule_completed'), is_schedule_completed),
  is_scheduled = COALESCE(sqlc.narg('is_scheduled'), is_scheduled),
  invitation_schedule_id = COALESCE(sqlc.narg('invitation_schedule_id'), invitation_schedule_id),
  invitation_schedule_gap = COALESCE(sqlc.narg('invitation_schedule_gap'), invitation_schedule_gap)
WHERE form_id = sqlc.arg('form_id')
RETURNING *;



-- name: GetFormByID :one
SELECT *
FROM forms
WHERE form_id = $1;

-- name: ListForms :many
SELECT
    *,
    COUNT(*) OVER() as total_count
FROM forms
WHERE
    created_by = sqlc.arg('created_by')
    AND (
        sqlc.narg('search')::text IS NULL
        OR form_title ILIKE '%' || sqlc.narg('search')::text || '%'
        OR form_description ILIKE '%' || sqlc.narg('search')::text || '%'
    )
    AND (sqlc.narg('form_status')::form_status IS NULL OR form_status = sqlc.narg('form_status')::form_status)
    AND (sqlc.narg('form_access')::form_access IS NULL OR form_access = sqlc.narg('form_access')::form_access)
ORDER BY
    CASE
        WHEN COALESCE(sqlc.narg('short_by')::text, '-updated') = 'updated'
            THEN form_updated_at
    END DESC,
    CASE
        WHEN COALESCE(sqlc.narg('short_by')::text, '-updated') = '-updated'
            THEN form_updated_at
    END ASC,
    CASE
        WHEN COALESCE(sqlc.narg('short_by')::text, '-updated') = 'title'
            THEN form_title
    END ASC,
    CASE
        WHEN COALESCE(sqlc.narg('short_by')::text, '-updated') = '-title'
            THEN form_title
    END DESC
LIMIT COALESCE(sqlc.narg('limit')::int, 10)
OFFSET COALESCE(sqlc.narg('offset')::int, 0);



-- name: DeleteForm :one
DELETE FROM forms
WHERE form_id = $1
RETURNING form_id, form_title, form_description, created_by, form_created_at, form_updated_at, form_status, form_access;


-- name: CreateFormField :one
INSERT INTO form_fields (form_id, field_label, field_type, is_required, ordering)
VALUES ($1, $2, $3, $4, $5)
RETURNING field_id, field_label, field_type, is_required, ordering, form_id;

-- name: UpdateFormField :one
UPDATE form_fields
SET
  field_label = COALESCE($2, field_label),
  field_type = COALESCE($3, field_type),
  is_required = COALESCE($4, is_required),
  ordering = COALESCE($5, ordering)
WHERE field_id = $1
RETURNING field_id, field_label, field_type, is_required, ordering, form_id;

-- name: DeleteFormField :one
DELETE FROM form_fields
WHERE field_id = $1
RETURNING field_id, field_label, field_type, is_required, ordering, form_id;

-- name: CreateFieldOption :one
INSERT INTO form_field_options (field_id, option_label, ordering, is_answer)
VALUES ($1, $2, $3, $4)
RETURNING option_id, field_id, option_label, ordering, is_answer;

-- name: UpdateFieldOption :one
UPDATE form_field_options
SET
  option_label = COALESCE($2, option_label),
  ordering = COALESCE($3, ordering),
  is_answer = COALESCE($4, is_answer)
WHERE option_id = $1
RETURNING option_id, field_id, option_label, ordering, is_answer;

-- name: DeleteFieldOption :one
DELETE FROM form_field_options
WHERE option_id = $1
RETURNING option_id, field_id, option_label, ordering, is_answer;


-- name: GetFormFieldsWithOptions :many
SELECT
    ff.field_id AS "fieldId",
    ff.field_label AS "fieldLabel",
    ff.field_type AS "fieldType",
    ff.is_required AS "isRequired",
    ff.ordering AS "ordering",
    ff.form_id AS "formId",
    COALESCE(
        JSON_AGG(
            JSONB_BUILD_OBJECT(
                'optionId', fo.option_id,
                'optionLabel', fo.option_label,
                'ordering', fo.ordering,
                'isAnswer', fo.is_answer,
                'fieldId', fo.field_id
            ) ORDER BY fo.ordering
        ) FILTER (WHERE fo.option_id IS NOT NULL),
        '[]'
    )::jsonb AS "options"
FROM form_fields ff
LEFT JOIN form_field_options fo ON ff.field_id = fo.field_id
WHERE ff.form_id = $1
GROUP BY ff.field_id
ORDER BY ff.ordering;

-- name: GetFormFieldOptionsMap :one
SELECT COALESCE(
    JSONB_OBJECT_AGG(
        fo.option_id::text,
        JSONB_BUILD_OBJECT(
            'optionId', fo.option_id,
            'optionLabel', fo.option_label,
            'ordering', fo.ordering,
            'fieldId', fo.field_id
        )
    ) FILTER (WHERE fo.option_id IS NOT NULL),
    '{}'::jsonb
) AS "fieldOptionsMap"
FROM form_field_options fo
JOIN form_fields ff ON ff.field_id = fo.field_id
WHERE ff.form_id = $1;
