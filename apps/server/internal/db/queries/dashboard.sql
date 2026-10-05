-- name: GetTotalForms :one
SELECT COUNT(*) FROM forms WHERE created_by = $1;

-- name: GetTotalSubmissions :one
SELECT COUNT(*) FROM form_submissions WHERE form_id IN (SELECT form_id FROM forms WHERE created_by = $1);

-- name: GetTotalActiveForms :one
SELECT COUNT(*) FROM forms WHERE created_by = $1 AND form_status = 'published';

-- name: GetTotalClosedForms :one
SELECT COUNT(*) FROM forms WHERE created_by = $1 AND form_status = 'closed';

-- name: GetTotalInvitations :one
SELECT COUNT(*) FROM invitations WHERE invited_by = $1;

-- name: GetFormSubmissionsByMonth :many
SELECT 
    TO_CHAR(submitted_at, 'YYYY-MM') AS month,
    COUNT(*) AS total_submissions
FROM form_submissions
WHERE form_id IN (SELECT form_id FROM forms WHERE created_by = $1)
GROUP BY month
ORDER BY month;
