ALTER TABLE note_read_models
    DROP COLUMN IF EXISTS notion_page_url;

ALTER TABLE notes
    DROP COLUMN IF EXISTS notion_page_id,
    DROP COLUMN IF EXISTS notion_page_url,
    DROP COLUMN IF EXISTS notion_synced_at;

ALTER TABLE templates
    DROP COLUMN IF EXISTS notion_parent_page_id;
