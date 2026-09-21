-- name: ListNotes :many
SELECT
    n.*,
    t.name AS template_name,
    a.first_name,
    a.last_name,
    a.thumbnail AS owner_thumbnail
FROM notes n
JOIN templates t ON t.id = n.template_id
JOIN accounts a ON a.id = n.owner_id
WHERE (NULLIF($1::text, '') IS NULL OR n.status = $1)
  AND ($2::uuid IS NULL OR n.template_id = $2)
  AND ($3::uuid IS NULL OR n.owner_id = $3)
  AND (NULLIF($4::text, '') IS NULL OR n.title ILIKE '%' || $4 || '%')
ORDER BY n.updated_at DESC;

-- name: GetNoteByID :one
SELECT
    n.*,
    t.name AS template_name,
    a.first_name,
    a.last_name,
    a.thumbnail AS owner_thumbnail
FROM notes n
JOIN templates t ON t.id = n.template_id
JOIN accounts a ON a.id = n.owner_id
WHERE n.id = $1;

-- name: CreateNote :one
INSERT INTO notes (title, template_id, owner_id, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateNote :one
UPDATE notes
SET
    title = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1;

-- name: UpdateNoteStatus :one
UPDATE notes
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListSectionsByNote :many
SELECT
    s.*,
    f.label,
    f."order",
    f.is_required
FROM sections s
JOIN fields f ON f.id = s.field_id
WHERE s.note_id = $1
ORDER BY f."order" ASC;

-- name: CreateSection :one
INSERT INTO sections (note_id, field_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSectionContent :one
UPDATE sections
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSectionsByNote :exec
DELETE FROM sections
WHERE note_id = $1;
