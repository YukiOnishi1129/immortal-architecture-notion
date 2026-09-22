-- Notion連携用のカラムを追加する。
-- 既存行を壊さないよう、すべて NULL 許容とする。

-- テンプレートごとの置き場所（Notionの親ページ）
ALTER TABLE templates
    ADD COLUMN notion_parent_page_id TEXT;

-- ノートごとに作られるNotionページ
ALTER TABLE notes
    ADD COLUMN notion_page_id TEXT,
    ADD COLUMN notion_page_url TEXT,
    ADD COLUMN notion_synced_at TIMESTAMPTZ;

-- 画面表示用（CQRSの読み取りモデル）。リンクのURLだけを持つ。
ALTER TABLE note_read_models
    ADD COLUMN notion_page_url TEXT;
