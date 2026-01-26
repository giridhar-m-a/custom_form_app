-- name: CreateManyInvitations :many
INSERT INTO invitations (
    form_id, 
    invited_email, 
    invited_name, 
    invited_by
)
SELECT
    @form_id::uuid,
    t1.email,
    t2.name,
    @invited_by::uuid
FROM
    UNNEST(@emails::text[]) WITH ORDINALITY AS t1(email, idx)
    JOIN UNNEST(@names::text[]) WITH ORDINALITY AS t2(name, idx) 
    USING (idx)
ON CONFLICT (invited_email, form_id) 
DO NOTHING
RETURNING invitation_id, invited_email, invited_name;

-- name: CreateInvitation :one
INSERT INTO invitations (
    form_id, 
    invited_email, 
    invited_name, 
    invited_by
)
VALUES (
    @form_id::uuid,
    @email::text,
    @name::text,
    @invited_by::uuid
)
ON CONFLICT (invited_email, form_id) 
DO NOTHING
RETURNING invitation_id, invited_email, invited_name;


-- name: UpdateInvitationStatus :one
UPDATE invitations
SET
    status = @status::invitation_status
WHERE
    invitation_id = @invitation_id::uuid
RETURNING invitation_id, invited_email, invited_name, status;

-- name: DeleteInvitation :exec
DELETE FROM invitations
WHERE invitation_id = @invitation_id::uuid
RETURNING invitation_id, invited_email, invited_name, status;

-- name: GetInvitationByFormId :many
SELECT * FROM invitations
WHERE form_id = @form_id::uuid
  -- Search filter (Email or Name)
  AND (
        sqlc.narg('search')::text IS NULL 
        OR invited_name ILIKE '%' || sqlc.narg('search')::text || '%'
        OR invited_email ILIKE '%' || sqlc.narg('search')::text || '%'
      )
  -- Status Inclusion filter
  AND (
        sqlc.narg('status')::invitation_status IS NULL 
        OR status = sqlc.narg('status')::invitation_status
      )
  -- Status Exclusion filter
  AND (
        sqlc.narg('exclude_status')::invitation_status IS NULL 
        OR status <> sqlc.narg('exclude_status')::invitation_status
      )
ORDER BY invited_at DESC
LIMIT COALESCE(sqlc.narg('limit_val')::int, 10)
OFFSET COALESCE(sqlc.narg('offset_val')::int, 0);



-- name: CountInvitationsByFormId :one
SELECT COUNT(*) AS total_records
FROM invitations
WHERE form_id = @form_id::uuid

  -- Search filter (Email or Name)
  AND (
        sqlc.narg('search')::text IS NULL 
        OR invited_name ILIKE '%' || sqlc.narg('search')::text || '%'
        OR invited_email ILIKE '%' || sqlc.narg('search')::text || '%'
      )

  -- Status Inclusion filter
  AND (
        sqlc.narg('status')::invitation_status IS NULL 
        OR status = sqlc.narg('status')::invitation_status
      )

  -- Status Exclusion filter
  AND (
        sqlc.narg('exclude_status')::invitation_status IS NULL 
        OR status <> sqlc.narg('exclude_status')::invitation_status
      );
