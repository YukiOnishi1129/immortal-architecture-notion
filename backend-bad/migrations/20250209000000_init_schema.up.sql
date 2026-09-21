-- Enable UUID generation for primary keys
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    provider TEXT NOT NULL,
    provider_account_id TEXT NOT NULL,
    thumbnail TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT provider_account_unique UNIQUE (provider, provider_account_id)
);

CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    owner_id UUID NOT NULL REFERENCES accounts(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    "order" INT NOT NULL CHECK ("order" > 0),
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fields_unique_order UNIQUE (template_id, "order")
);

CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    template_id UUID NOT NULL REFERENCES templates(id),
    owner_id UUID NOT NULL REFERENCES accounts(id),
    status TEXT NOT NULL CHECK (status IN ('Draft', 'Publish')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    note_id UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    field_id UUID NOT NULL REFERENCES fields(id),
    content TEXT NOT NULL DEFAULT '',
    CONSTRAINT sections_unique_field UNIQUE (note_id, field_id)
);

CREATE INDEX idx_templates_owner_id ON templates(owner_id);
CREATE INDEX idx_fields_template_id ON fields(template_id);
CREATE INDEX idx_notes_template_id ON notes(template_id);
CREATE INDEX idx_notes_owner_id ON notes(owner_id);
CREATE INDEX idx_notes_updated_at ON notes(updated_at DESC);
CREATE INDEX idx_sections_note_id ON sections(note_id);
CREATE INDEX idx_sections_field_id ON sections(field_id);
