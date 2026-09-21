DROP INDEX IF EXISTS idx_sections_field_id;
DROP INDEX IF EXISTS idx_sections_note_id;
DROP INDEX IF EXISTS idx_notes_updated_at;
DROP INDEX IF EXISTS idx_notes_owner_id;
DROP INDEX IF EXISTS idx_notes_template_id;
DROP INDEX IF EXISTS idx_fields_template_id;
DROP INDEX IF EXISTS idx_templates_owner_id;

DROP TABLE IF EXISTS sections;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS fields;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS accounts;
