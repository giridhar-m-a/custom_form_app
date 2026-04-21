-- name: CreateSubmission :one
INSERT INTO form_submissions (
    form_id,
    respondent_id
) VALUES (
    sqlc.arg(form_id)::uuid,
    sqlc.narg(respondent_id)::uuid
) RETURNING *;

-- name: CreateResponse :one
INSERT INTO form_responses (
    submission_id,
    form_field_id,
    response_text
) VALUES (
    @submission_id::uuid,
    @form_field_id::uuid,
    sqlc.narg('response_text')::text
) RETURNING *;

-- name: CreateResponseOption :one
INSERT INTO response_options (
    response_id,
    form_option_id
) VALUES (
    @response_id::uuid,
    @form_option_id::uuid
) RETURNING *;

-- name: CreateResponseFiles :one
INSERT INTO form_response_files (
    response_id,
    file_name,
    file_path,
    file_size,
    file_type,
    form_id
) VALUES (
    @response_id::uuid,
    @file_name::varchar,
    @file_path::varchar,
    @file_size::bigint,
    @file_type::varchar,
    @form_id::uuid
) RETURNING *;


-- name: GetSubmissions :many
WITH filtered AS (
    SELECT 
        fs.*,
        i.invited_email,
        i.invited_name,
        i.status         AS invitation_status,
        i.invited_at,
        i.opened_at      AS invitation_opened_at
    FROM form_submissions fs
    LEFT JOIN invitations i 
        ON fs.respondent_id = i.invitation_id
    WHERE 
        fs.form_id = sqlc.arg(form_id)::uuid
        AND (
            sqlc.narg(search)::text IS NULL
            OR i.invited_email ILIKE CONCAT('%', sqlc.narg(search), '%')
            OR i.invited_name ILIKE CONCAT('%', sqlc.narg(search), '%')
        )
)
SELECT * FROM filtered
ORDER BY submitted_at DESC
LIMIT sqlc.arg(limit_count)::int
OFFSET sqlc.arg(offset_count)::int;

-- name: GetSubmissionCount :one
WITH filtered AS (
    SELECT fs.submission_id
    FROM form_submissions fs
    LEFT JOIN invitations i 
        ON fs.respondent_id = i.invitation_id
    WHERE 
        fs.form_id = sqlc.arg(form_id)::uuid
        AND (
            sqlc.narg(search)::text IS NULL
            OR i.invited_email ILIKE CONCAT('%', sqlc.narg(search), '%')
            OR i.invited_name ILIKE CONCAT('%', sqlc.narg(search), '%')
        )
)
SELECT COUNT(*) FROM filtered;


-- name: GetSubmissionById :one
SELECT 
    jsonb_build_object(
        'submission_id', fs.submission_id,
        'form_id', fs.form_id,
        'submitted_at', fs.submitted_at,
        'respondent_id', fs.respondent_id,
        'responses', COALESCE(responses.responses, '[]'::jsonb)
    ) AS submission
FROM form_submissions fs
LEFT JOIN LATERAL (
    SELECT jsonb_agg(
        jsonb_build_object(
            'response_id', fr.response_id,
            'form_field_id', fr.form_field_id,
            'response_text', fr.response_text,

            -- options (multi-select)
            'form_field_options', (
                SELECT COALESCE(jsonb_agg(
                    jsonb_build_object(
                        'id', ro.id,
                        'form_option_id', ro.form_option_id
                    )
                ), '[]'::jsonb)
                FROM response_options ro
                WHERE ro.response_id = fr.response_id
            ),

            -- files
            'form_response_files', (
                SELECT COALESCE(jsonb_agg(
                    jsonb_build_object(
                        'response_file_id', frf.response_file_id,
                        'file_name', frf.file_name,
                        'file_path', frf.file_path,
                        'file_size', frf.file_size,
                        'file_type', frf.file_type,
                        'file_uploaded_at', frf.file_uploaded_at
                    )
                ), '[]'::jsonb)
                FROM form_response_files frf
                WHERE frf.response_id = fr.response_id
            )
        )
    ) AS responses
    FROM form_responses fr
    WHERE fr.submission_id = fs.submission_id
) responses ON TRUE
WHERE fs.submission_id = $1;