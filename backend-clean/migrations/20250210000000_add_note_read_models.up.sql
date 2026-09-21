-- Read model table for CQRS: denormalized note data for fast reads without JOINs.
CREATE TABLE note_read_models (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('Draft', 'Publish')),
    template_id UUID NOT NULL,
    template_name TEXT NOT NULL,
    owner_id UUID NOT NULL,
    owner_first_name TEXT NOT NULL,
    owner_last_name TEXT NOT NULL,
    owner_thumbnail TEXT,
    sections_json JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_note_read_models_status ON note_read_models(status);
CREATE INDEX idx_note_read_models_owner_id ON note_read_models(owner_id);
CREATE INDEX idx_note_read_models_template_id ON note_read_models(template_id);
CREATE INDEX idx_note_read_models_updated_at ON note_read_models(updated_at DESC);
