-- name: ListTemplates :many
SELECT
    t.*,
    EXISTS (
        SELECT 1
        FROM notes n
        WHERE n.template_id = t.id
        LIMIT 1
    ) AS is_used
FROM templates t
WHERE ($1::uuid IS NULL OR t.owner_id = $1)
  AND ($2::text IS NULL OR t.name ILIKE '%' || $2 || '%')
ORDER BY t.updated_at DESC;

-- name: GetTemplateByID :one
SELECT
    t.*,
    EXISTS (
        SELECT 1
        FROM notes n
        WHERE n.template_id = t.id
        LIMIT 1
    ) AS is_used
FROM templates t
WHERE t.id = $1;

-- name: CreateTemplate :one
INSERT INTO templates (name, owner_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTemplate :one
UPDATE templates
SET
    name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTemplate :exec
DELETE FROM templates
WHERE id = $1;

-- name: CheckTemplateInUse :one
SELECT EXISTS (
    SELECT 1 FROM notes WHERE template_id = $1
) AS is_used;

-- name: ListFieldsByTemplate :many
SELECT *
FROM fields
WHERE template_id = $1
ORDER BY "order" ASC;

-- name: CreateField :one
INSERT INTO fields (template_id, label, "order", is_required)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateField :one
UPDATE fields
SET
    label = $2,
    "order" = $3,
    is_required = $4
WHERE id = $1
RETURNING *;

-- name: DeleteFieldsByTemplate :exec
DELETE FROM fields
WHERE template_id = $1;

-- name: DeleteField :exec
DELETE FROM fields
WHERE id = $1;
