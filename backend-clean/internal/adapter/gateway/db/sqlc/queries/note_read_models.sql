-- name: ListNoteReadModels :many
SELECT *
FROM note_read_models
WHERE (NULLIF($1::text, '') IS NULL OR status = $1)
  AND ($2::uuid IS NULL OR template_id = $2)
  AND ($3::uuid IS NULL OR owner_id = $3)
  AND (NULLIF($4::text, '') IS NULL OR title ILIKE '%' || $4 || '%')
ORDER BY updated_at DESC;

-- name: GetNoteReadModel :one
SELECT *
FROM note_read_models
WHERE id = $1;

-- name: UpsertNoteReadModel :exec
INSERT INTO note_read_models (
    id, title, status, template_id, template_name,
    owner_id, owner_first_name, owner_last_name, owner_thumbnail,
    sections_json, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    status = EXCLUDED.status,
    template_id = EXCLUDED.template_id,
    template_name = EXCLUDED.template_name,
    owner_id = EXCLUDED.owner_id,
    owner_first_name = EXCLUDED.owner_first_name,
    owner_last_name = EXCLUDED.owner_last_name,
    owner_thumbnail = EXCLUDED.owner_thumbnail,
    sections_json = EXCLUDED.sections_json,
    updated_at = EXCLUDED.updated_at;

-- name: DeleteNoteReadModel :exec
DELETE FROM note_read_models
WHERE id = $1;
