-- Seed note_read_models from existing data (one-time sync for CQRS migration).
INSERT INTO note_read_models (
    id, title, status, template_id, template_name,
    owner_id, owner_first_name, owner_last_name, owner_thumbnail,
    sections_json, created_at, updated_at
)
SELECT
    n.id,
    n.title,
    n.status,
    n.template_id,
    t.name,
    n.owner_id,
    a.first_name,
    a.last_name,
    a.thumbnail,
    COALESCE(
        (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'id', s.id,
                    'field_id', s.field_id,
                    'field_label', f.label,
                    'field_order', f."order",
                    'is_required', f.is_required,
                    'content', s.content
                ) ORDER BY f."order"
            )
            FROM sections s
            JOIN fields f ON f.id = s.field_id
            WHERE s.note_id = n.id
        ),
        '[]'::jsonb
    ),
    n.created_at,
    n.updated_at
FROM notes n
JOIN templates t ON t.id = n.template_id
JOIN accounts a ON a.id = n.owner_id
ON CONFLICT (id) DO NOTHING;
